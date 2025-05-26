package _handlers

import (
	"database/sql"
	"strings"

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

		// Write CSV header
		for i, col := range columns {
			if i > 0 {
				c.Writer.WriteString(",")
			}
			c.Writer.WriteString(col)
		}
		c.Writer.WriteString("\n")

		// Query complète avec toutes les colonnes
		rows, err := db.Query(`
			SELECT `+strings.Join(columns, ",")+`
			FROM `+tableName+`
			WHERE `+columnName+` ILIKE $1
		`, "%"+searchValue+"%")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		// Écrire toutes les colonnes, pas juste une
		for rows.Next() {
			values := make([]sql.NullString, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}

			// Write CSV line
			for i, val := range values {
				if i > 0 {
					c.Writer.WriteString(",")
				}
				if val.Valid {
					// Escape quotes in CSV
					escaped := strings.ReplaceAll(val.String, `"`, `""`)
					if strings.Contains(escaped, ",") || strings.Contains(escaped, "\n") || strings.Contains(escaped, `"`) {
						c.Writer.WriteString(`"` + escaped + `"`)
					} else {
						c.Writer.WriteString(escaped)
					}
				}
			}
			c.Writer.WriteString("\n")
		}
	}
}
