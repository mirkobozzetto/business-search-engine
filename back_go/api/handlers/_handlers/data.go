package _handlers

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PreviewTable(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		limitStr := c.DefaultQuery("limit", "5")
		limit, _ := strconv.Atoi(limitStr)

		rows, err := db.Query(`SELECT * FROM `+tableName+` LIMIT $1`, limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		var data []gin.H

		for rows.Next() {
			values := make([]sql.NullString, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}

			row := gin.H{}
			for i, col := range columns {
				if values[i].Valid {
					row[col] = values[i].String
				} else {
					row[col] = nil
				}
			}
			data = append(data, row)
		}

		c.JSON(200, gin.H{
			"table":   tableName,
			"columns": columns,
			"data":    data,
		})
	}
}

func GetColumnValues(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		limitStr := c.DefaultQuery("limit", "20")
		limit, _ := strconv.Atoi(limitStr)

		rows, err := db.Query(`
			SELECT `+columnName+`, count(*) as count
			FROM `+tableName+`
			WHERE `+columnName+` IS NOT NULL AND `+columnName+` != ''
			GROUP BY `+columnName+`
			ORDER BY count DESC
			LIMIT $1
		`, limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var values []gin.H
		for rows.Next() {
			var value string
			var count int64
			if err := rows.Scan(&value, &count); err != nil {
				continue
			}
			values = append(values, gin.H{
				"value": value,
				"count": count,
			})
		}

		c.JSON(200, gin.H{
			"table":  tableName,
			"column": columnName,
			"values": values,
		})
	}
}
