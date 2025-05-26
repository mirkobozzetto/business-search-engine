package handlers

import (
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ListTables(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT table_name,
				   (SELECT count(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
			FROM information_schema.tables t
			WHERE table_schema = 'public'
			ORDER BY table_name
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(500, models.Error("failed to get tables: "+err.Error()))
			return
		}
		defer rows.Close()

		var result []models.Table
		for rows.Next() {
			var tableName string
			var columnCount int
			if err := rows.Scan(&tableName, &columnCount); err != nil {
				continue
			}

			// Get row count
			var rowCount int64
			rowQuery := fmt.Sprintf("SELECT count(*) FROM %s", tableName)
			db.QueryRow(rowQuery).Scan(&rowCount)

			result = append(result, models.Table{
				Name:    tableName,
				Rows:    rowCount,
				Columns: columnCount,
			})
		}

		c.JSON(200, models.Success(result))
	}
}

func GetTableInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("name")

		// Validate table name
		if err := helpers.ValidateTableName(tableName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			return
		}

		// Get column count
		colQuery := `SELECT count(*) FROM information_schema.columns WHERE table_name = $1`
		var colCount int
		db.QueryRow(colQuery, tableName).Scan(&colCount)

		// Get row count
		rowQuery := fmt.Sprintf("SELECT count(*) FROM %s", tableName)
		var rowCount int64
		db.QueryRow(rowQuery).Scan(&rowCount)

		// Get column names
		fieldsQuery := `SELECT column_name FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position`
		rows, err := db.Query(fieldsQuery, tableName)
		if err != nil {
			c.JSON(500, models.Error("failed to get fields: "+err.Error()))
			return
		}
		defer rows.Close()

		var fields []string
		for rows.Next() {
			var field string
			if err := rows.Scan(&field); err != nil {
				continue
			}
			fields = append(fields, field)
		}

		result := models.TableInfo{
			Table:   tableName,
			Rows:    rowCount,
			Columns: colCount,
			Fields:  fields,
		}

		c.JSON(200, models.Success(result))
	}
}

func GetTableColumns(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("name")

		// Validate table name
		if err := helpers.ValidateTableName(tableName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			return
		}

		query := `
			SELECT column_name, data_type, is_nullable
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position
		`

		rows, err := db.Query(query, tableName)
		if err != nil {
			c.JSON(500, models.Error("failed to get columns: "+err.Error()))
			return
		}
		defer rows.Close()

		var result []models.ColumnInfo
		for rows.Next() {
			var name, dataType, nullable string
			if err := rows.Scan(&name, &dataType, &nullable); err != nil {
				continue
			}

			result = append(result, models.ColumnInfo{
				Name:     name,
				Type:     dataType,
				Nullable: nullable == "YES",
			})
		}

		c.JSON(200, models.Success(result))
	}
}

func GetCompleteStructure(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get all tables with their info
		tablesQuery := `
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'public'
			ORDER BY table_name
		`

		rows, err := db.Query(tablesQuery)
		if err != nil {
			c.JSON(500, models.Error("failed to get tables: "+err.Error()))
			return
		}
		defer rows.Close()

		var result []models.TableStructure
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				continue
			}

			// Get row count
			var rowCount int64
			rowQuery := fmt.Sprintf("SELECT count(*) FROM %s", tableName)
			db.QueryRow(rowQuery).Scan(&rowCount)

			// Get columns
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
				var name, dataType, nullable string
				if err := colRows.Scan(&name, &dataType, &nullable); err != nil {
					continue
				}
				columns = append(columns, models.ColumnInfo{
					Name:     name,
					Type:     dataType,
					Nullable: nullable == "YES",
				})
			}
			colRows.Close()

			result = append(result, models.TableStructure{
				Name:    tableName,
				Rows:    rowCount,
				Columns: columns,
			})
		}

		c.JSON(200, models.Success(result))
	}
}
