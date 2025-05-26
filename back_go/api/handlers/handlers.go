package handlers

import (
	"csv-importer/api/handlers/_handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ListTables(db *sql.DB) gin.HandlerFunc {
	return _handlers.ListTables(db)
}

func GetTableInfo(db *sql.DB) gin.HandlerFunc {
	return _handlers.GetTableInfo(db)
}

func GetTableColumns(db *sql.DB) gin.HandlerFunc {
	return _handlers.GetTableColumns(db)
}

func PreviewTable(db *sql.DB) gin.HandlerFunc {
	return _handlers.PreviewTable(db)
}

func GetColumnValues(db *sql.DB) gin.HandlerFunc {
	return _handlers.GetColumnValues(db)
}

func SearchTable(db *sql.DB) gin.HandlerFunc {
	return _handlers.SearchTable(db)
}

func CountRows(db *sql.DB) gin.HandlerFunc {
	return _handlers.CountRows(db)
}

func ExportData(db *sql.DB) gin.HandlerFunc {
	return _handlers.ExportData(db)
}

func GetCompleteStructure(db *sql.DB) gin.HandlerFunc {
	return _handlers.GetCompleteStructure(db)
}
