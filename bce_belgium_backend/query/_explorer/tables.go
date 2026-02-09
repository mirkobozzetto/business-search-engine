package _explorer

import (
	"database/sql"
	"fmt"
	"strings"
)

func ListTables(db *sql.DB) error {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying tables: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n📊 TABLES:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		fmt.Printf("   • %s\n", tableName)
	}

	return nil
}

func ShowTableInfo(db *sql.DB, tableName string) error {
	// Get column count
	colQuery := `SELECT count(*) FROM information_schema.columns WHERE table_name = $1`
	var colCount int
	db.QueryRow(colQuery, tableName).Scan(&colCount)

	// Get row count
	rowQuery := fmt.Sprintf("SELECT count(*) FROM %s", tableName)
	var rowCount int64
	db.QueryRow(rowQuery).Scan(&rowCount)

	fmt.Printf("\n🔍 TABLE: %s\n", strings.ToUpper(tableName))
	fmt.Printf("   📊 Rows: %d\n", rowCount)
	fmt.Printf("   📋 Columns: %d\n", colCount)

	return nil
}

func ShowStats(db *sql.DB) error {
	query := `
		SELECT tablename, n_live_tup
		FROM pg_stat_user_tables
		ORDER BY n_live_tup DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error getting stats: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n🚀 STATS:")
	fmt.Printf("%-20s %s\n", "TABLE", "ROWS")
	fmt.Println(strings.Repeat("-", 30))

	var totalRows int64
	for rows.Next() {
		var table string
		var liveRows int64

		if err := rows.Scan(&table, &liveRows); err != nil {
			continue
		}

		fmt.Printf("%-20s %d\n", table, liveRows)
		totalRows += liveRows
	}

	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("TOTAL: %d rows\n", totalRows)

	return nil
}
