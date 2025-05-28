package handlers

import (
	"csv-importer/api/helpers"
	helperutils "csv-importer/api/helpers/utils"
	"csv-importer/api/models"
	"csv-importer/api/utils"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchTable(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")
		limitStr := c.DefaultQuery("limit", "50")

		if searchValue == "" {
			c.JSON(400, models.Error("search query 'q' is required"))
			return
		}

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

		// Safe search query
		query := fmt.Sprintf(`
			SELECT DISTINCT %s
			FROM %s
			WHERE %s ILIKE $1
			ORDER BY %s
			LIMIT $2
		`, columnName, tableName, columnName, columnName)

		rows, err := db.Query(query, "%"+searchValue+"%", limit)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}
		defer rows.Close()

		var results []string
		for rows.Next() {
			var value sql.NullString
			if err := rows.Scan(&value); err != nil {
				continue
			}
			if value.Valid {
				results = append(results, value.String)
			}
		}

		result := models.SearchResult{
			Table:   tableName,
			Column:  columnName,
			Query:   searchValue,
			Results: results,
			Meta: models.Meta{
				Count: len(results),
				Limit: limit,
			},
		}

		c.JSON(200, models.Success(result))
	}
}

func CountRows(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")

		if searchValue == "" {
			c.JSON(400, models.Error("search query 'q' is required"))
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

		// Safe count query
		query := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM %s
			WHERE %s ILIKE $1
		`, tableName, columnName)

		var count int64
		err := db.QueryRow(query, "%"+searchValue+"%").Scan(&count)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}

		result := models.CountResult{
			Table:  tableName,
			Column: columnName,
			Query:  searchValue,
			Count:  count,
		}

		c.JSON(200, models.Success(result))
	}
}

func SearchNaceCode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		searchValue := c.Query("q")
		limitStr := c.Query("limit")

		limit := utils.ParseOptionalLimit(limitStr, 0)

		query, args := utils.BuildNaceCodeQuery(searchValue, limit)

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}
		defer rows.Close()

		columns := []string{"nacecode", "activités", "libellé_fr", "omschrijving_nl"}
		data, err := helperutils.ScanRowsToMaps(rows, columns)
		if err != nil {
			c.JSON(500, models.Error("scan error: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(map[string]any{
			"query":   searchValue,
			"results": data,
			"meta": models.Meta{
				Count: len(data),
				Limit: limit,
				Total: len(data),
			},
		}))
	}
}
