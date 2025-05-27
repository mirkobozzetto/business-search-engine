package handlers

import (
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PreviewTable(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		limitStr := c.DefaultQuery("limit", "5")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		// Validate table name
		if err := helpers.ValidateTableName(tableName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			return
		}

		// Safe query using validated table name with LIMIT
		builder := &helpers.QueryBuilder{}
		builder.SetLimit(limit)
		rows, err := helpers.SafeQueryWithBuilder(db, tableName, nil, builder)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			c.JSON(500, models.Error("failed to get columns"))
			return
		}

		var data []map[string]any
		rowCount := 0

		for rows.Next() {
			values := make([]sql.NullString, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}

			row := make(map[string]any)
			for i, col := range columns {
				if values[i].Valid {
					row[col] = values[i].String
				} else {
					row[col] = nil
				}
			}
			data = append(data, row)
			rowCount++
		}

		result := models.PreviewData{
			Table:   tableName,
			Columns: columns,
			Data:    data,
			Meta: models.Meta{
				Count: rowCount,
				Limit: limit,
			},
		}

		c.JSON(200, models.Success(result))
	}
}

func GetColumnValues(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		limitStr := c.DefaultQuery("limit", "20")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		// Validate table and column
		if err := helpers.ValidateTableName(tableName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			return
		}
		if err := helpers.ValidateColumnExists(db, tableName, columnName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			return
		}

		// Safe query using QueryBuilder
		builder := &helpers.QueryBuilder{}
		builder.SetLimit(limit)
		// Note: GROUP BY and ORDER BY are not supported by QueryBuilder
		// Direct usage but secure because tableName and columnName are validated
		query := fmt.Sprintf(`
			SELECT %s, count(*) as count
			FROM %s
			WHERE %s IS NOT NULL AND %s != ''
			GROUP BY %s
			ORDER BY count DESC
			LIMIT $1
		`, columnName, tableName, columnName, columnName, columnName)
		rows, err := db.Query(query, limit)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}
		defer rows.Close()

		var values []models.ColumnValue
		for rows.Next() {
			var value string
			var count int64
			if err := rows.Scan(&value, &count); err != nil {
				continue
			}
			values = append(values, models.ColumnValue{
				Value: value,
				Count: count,
			})
		}

		result := models.ColumnValues{
			Table:  tableName,
			Column: columnName,
			Values: values,
			Meta: models.Meta{
				Count: len(values),
				Limit: limit,
			},
		}

		c.JSON(200, models.Success(result))
	}
}
