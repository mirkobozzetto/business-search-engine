package query

import (
	"csv-importer/query/_explorer"
	"database/sql"
)

// Tables functions
func ListTables(db *sql.DB) error {
	return _explorer.ListTables(db)
}

func ShowTableInfo(db *sql.DB, tableName string) error {
	return _explorer.ShowTableInfo(db, tableName)
}

func ShowStats(db *sql.DB) error {
	return _explorer.ShowStats(db)
}

// Columns functions
func ShowColumns(db *sql.DB, tableName string) error {
	return _explorer.ShowColumns(db, tableName)
}

func ShowColumnValues(db *sql.DB, tableName, columnName string, limit int) error {
	return _explorer.ShowColumnValues(db, tableName, columnName, limit)
}

// Data functions
func PreviewTable(db *sql.DB, tableName string, limit int) error {
	return _explorer.PreviewTable(db, tableName, limit)
}

func SearchTable(db *sql.DB, tableName, columnName, searchValue string, limit int) error {
	return _explorer.SearchTable(db, tableName, columnName, searchValue, limit)
}

// Analysis functions
func CountRows(db *sql.DB, tableName, columnName, searchValue string) error {
	return _explorer.CountRows(db, tableName, columnName, searchValue)
}

func SampleRows(db *sql.DB, tableName, columnName, searchValue string, limit int) error {
	return _explorer.SampleRows(db, tableName, columnName, searchValue, limit)
}

func ExportToCSV(db *sql.DB, tableName, columnName, searchValue, filename string) error {
	return _explorer.ExportToCSV(db, tableName, columnName, searchValue, filename)
}
