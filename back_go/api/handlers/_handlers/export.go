package _handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ExportData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")

		if searchValue == "" {
			c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", `attachment; filename="`+tableName+`_export.csv"`)

		// Get columns
		colRows, err := db.Query(`
			SELECT column_name
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position
		`, tableName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer colRows.Close()

		var columns []string
		for colRows.Next() {
			var col string
			if err := colRows.Scan(&col); err != nil {
				continue
			}
			columns = append(columns, col)
		}

		for i, col := range columns {
			if i > 0 {
				c.Writer.WriteString(",")
			}
			c.Writer.WriteString(col)
		}
		c.Writer.WriteString("\n")

		dataRows, err := db.Query(`
			SELECT `+columnName+`
			FROM `+tableName+`
			WHERE `+columnName+` ILIKE $1
		`, "%"+searchValue+"%")
		if err != nil {
			return
		}
		defer dataRows.Close()

		for dataRows.Next() {
			var value string
			if err := dataRows.Scan(&value); err != nil {
				continue
			}
			c.Writer.WriteString(value + "\n")
		}
	}
}
