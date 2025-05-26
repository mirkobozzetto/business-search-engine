package _handlers

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchTable(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")
		limitStr := c.DefaultQuery("limit", "10")
		limit, _ := strconv.Atoi(limitStr)

		if searchValue == "" {
			c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
			return
		}

		rows, err := db.Query(`
			SELECT `+columnName+`
			FROM `+tableName+`
			WHERE `+columnName+` ILIKE $1
			LIMIT $2
		`, "%"+searchValue+"%", limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var results []string
		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err != nil {
				continue
			}
			results = append(results, value)
		}

		c.JSON(200, gin.H{
			"table":   tableName,
			"column":  columnName,
			"query":   searchValue,
			"results": results,
			"count":   len(results),
		})
	}
}

func CountRows(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")

		if searchValue == "" {
			c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
			return
		}

		var count int64
		err := db.QueryRow(`
			SELECT count(*)
			FROM `+tableName+`
			WHERE `+columnName+` ILIKE $1
		`, "%"+searchValue+"%").Scan(&count)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"table":  tableName,
			"column": columnName,
			"query":  searchValue,
			"count":  count,
		})
	}
}
