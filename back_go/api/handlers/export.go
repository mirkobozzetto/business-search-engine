package handlers

import (
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExportData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Query("column")
		searchValue := c.Query("search")
		limitStr := c.DefaultQuery("limit", "10000")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100000 {
			c.JSON(400, models.Error("invalid limit parameter (max 100000)"))
			return
		}

		// Validate table name
		if err := helpers.ValidateTableName(tableName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			return
		}

		// Get all columns for the table
		columnsQuery := `
			SELECT column_name
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position
		`
		rows, err := db.Query(columnsQuery, tableName)
		if err != nil {
			c.JSON(500, models.Error("failed to get table columns: "+err.Error()))
			return
		}
		defer rows.Close()

		var columns []string
		for rows.Next() {
			var column string
			if err := rows.Scan(&column); err != nil {
				continue
			}
			columns = append(columns, column)
		}

		if len(columns) == 0 {
			c.JSON(404, models.Error("table not found or has no columns"))
			return
		}

		// Build query
		var query string
		var args []any

		if columnName != "" && searchValue != "" {
			// Validate the search column exists
			if err := helpers.ValidateColumnExists(db, tableName, columnName); err != nil {
				c.JSON(400, models.Error(err.Error()))
				return
			}

			query = fmt.Sprintf(`
				SELECT %s
				FROM %s
				WHERE %s ILIKE $1
				LIMIT $2
			`, strings.Join(columns, ","), tableName, columnName)
			args = []any{"%" + searchValue + "%", limit}
		} else {
			query = fmt.Sprintf(`
				SELECT %s
				FROM %s
				LIMIT $1
			`, strings.Join(columns, ","), tableName)
			args = []any{limit}
		}

		// Execute query
		dataRows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}
		defer dataRows.Close()

		// Set CSV headers
		filename := fmt.Sprintf("%s_export.csv", tableName)
		if searchValue != "" {
			filename = fmt.Sprintf("%s_%s_export.csv", tableName, searchValue)
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		// Create CSV writer
		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		// Write header
		if err := writer.Write(columns); err != nil {
			c.JSON(500, models.Error("failed to write CSV header"))
			return
		}

		// Write data
		rowCount := 0
		for dataRows.Next() {
			values := make([]sql.NullString, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := dataRows.Scan(valuePtrs...); err != nil {
				continue
			}

			record := make([]string, len(columns))
			for i, val := range values {
				if val.Valid {
					record[i] = val.String
				} else {
					record[i] = ""
				}
			}

			if err := writer.Write(record); err != nil {
				break
			}
			rowCount++
		}

		// If no data was written, we need to handle this differently
		// since we already set headers for CSV download
		if rowCount == 0 {
			// Write an empty row to indicate no data
			writer.Write([]string{"No data found"})
		}

		c.Status(http.StatusOK)
	}
}
