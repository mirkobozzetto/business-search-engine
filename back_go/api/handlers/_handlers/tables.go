package _handlers

import (
	"csv-importer/api/models"
	"database/sql"

	"github.com/gin-gonic/gin"
)


func ListTables(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT table_name,
				   (SELECT count(*) FROM information_schema.columns WHERE table_name = t.table_name) as columns,
				   pg_stat_get_live_tuples(c.oid) as rows
			FROM information_schema.tables t
			LEFT JOIN pg_class c ON c.relname = t.table_name
			WHERE table_schema = 'public'
			ORDER BY table_name
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var tables []gin.H
		for rows.Next() {
			var name string
			var columns, tableRows sql.NullInt64

			if err := rows.Scan(&name, &columns, &tableRows); err != nil {
				continue
			}

			tables = append(tables, gin.H{
				"name":    name,
				"columns": columns.Int64,
				"rows":    tableRows.Int64,
			})
		}

		c.JSON(200, gin.H{"tables": tables})
	}
}

func GetTableInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("name")

		// Get column count
		var colCount int
		db.QueryRow(`SELECT count(*) FROM information_schema.columns WHERE table_name = $1`, tableName).Scan(&colCount)

		// Get row count
		var rowCount int64
		db.QueryRow(`SELECT count(*) FROM `+tableName).Scan(&rowCount)

		c.JSON(200, gin.H{
			"table":   tableName,
			"rows":    rowCount,
			"columns": colCount,
		})
	}
}

func GetTableColumns(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("name")

		rows, err := db.Query(`
			SELECT column_name
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position
		`, tableName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var columns []string
		for rows.Next() {
			var col string
			if err := rows.Scan(&col); err != nil {
				continue
			}
			columns = append(columns, col)
		}

		c.JSON(200, gin.H{"columns": columns})
	}
}


func GetCompleteStructure(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT t.table_name,
				   (SELECT count(*) FROM information_schema.columns WHERE table_name = t.table_name) as columns,
				   pg_stat_get_live_tuples(cl.oid) as rows
			FROM information_schema.tables t
			LEFT JOIN pg_class cl ON cl.relname = t.table_name
			WHERE t.table_schema = 'public'
			ORDER BY t.table_name
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(500, models.Error("error getting tables: "+err.Error()))
			return
		}
		defer rows.Close()

		var result []models.TableStructure

		for rows.Next() {
			var tableName string
			var columnCount, rowCount sql.NullInt64

			if err := rows.Scan(&tableName, &columnCount, &rowCount); err != nil {
				continue
			}

			colQuery := `
				SELECT column_name, data_type, is_nullable
				FROM information_schema.columns
				WHERE table_name = $1
				ORDER BY ordinal_position
			`

			colRows, err := db.Query(colQuery, tableName)
			if err != nil {
				continue
			}

			var columns []models.ColumnInfo
			for colRows.Next() {
				var colName, dataType, isNullable string
				if err := colRows.Scan(&colName, &dataType, &isNullable); err != nil {
					continue
				}

				columns = append(columns, models.ColumnInfo{
					Name:     colName,
					Type:     dataType,
					Nullable: isNullable == "YES",
				})
			}
			colRows.Close()

			result = append(result, models.TableStructure{
				Name:    tableName,
				Rows:    rowCount.Int64,
				Columns: columns,
			})
		}

		c.JSON(200, models.Success(result))
	}
}
