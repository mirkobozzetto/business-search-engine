package _explorer

import (
	"database/sql"
	"fmt"
	"strings"
)

func ShowColumns(db *sql.DB, tableName string) error {
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return fmt.Errorf("error getting columns: %v", err)
	}
	defer rows.Close()

	fmt.Printf("\n📋 COLUMNS in %s:\n", strings.ToUpper(tableName))
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			continue
		}
		fmt.Printf("   • %s\n", columnName)
	}

	return nil
}

func ShowColumnValues(db *sql.DB, tableName, columnName string, limit int) error {
	if limit <= 0 {
		limit = 20
	}

	fmt.Printf("🔍 Checking table '%s' column '%s'...\n", tableName, columnName)

	// First check if table exists
	tableCheck := `SELECT count(*) FROM information_schema.tables WHERE table_name = $1`
	var tableCount int
	if err := db.QueryRow(tableCheck, tableName).Scan(&tableCount); err != nil {
		return fmt.Errorf("error checking table: %v", err)
	}
	if tableCount == 0 {
		return fmt.Errorf("table '%s' not found", tableName)
	}

	// Check if column exists
	colCheck := `SELECT count(*) FROM information_schema.columns WHERE table_name = $1 AND column_name = $2`
	var colCount int
	if err := db.QueryRow(colCheck, tableName, columnName).Scan(&colCount); err != nil {
		return fmt.Errorf("error checking column: %v", err)
	}
	if colCount == 0 {
		return fmt.Errorf("column '%s' not found in table '%s'", columnName, tableName)
	}

	fmt.Printf("✅ Table and column found, getting values...\n")

	query := fmt.Sprintf(`
		SELECT %s, count(*) as count
		FROM %s
		WHERE %s IS NOT NULL AND %s != ''
		GROUP BY %s
		ORDER BY count DESC
		LIMIT %d
	`, columnName, tableName, columnName, columnName, columnName, limit)

	fmt.Printf("📝 SQL: %s\n", query)

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying column values: %v", err)
	}
	defer rows.Close()

	fmt.Printf("\n🔍 VALUES: %s.%s (top %d)\n", strings.ToUpper(tableName), strings.ToUpper(columnName), limit)
	fmt.Printf("%-30s %s\n", "VALUE", "COUNT")
	fmt.Println(strings.Repeat("-", 40))

	resultCount := 0
	for rows.Next() {
		var value string
		var count int64

		if err := rows.Scan(&value, &count); err != nil {
			fmt.Printf("❌ Error scanning row: %v\n", err)
			continue
		}

		if len(value) > 27 {
			value = value[:24] + "..."
		}

		fmt.Printf("%-30s %d\n", value, count)
		resultCount++
	}

	if resultCount == 0 {
		fmt.Println("❌ No data found - table might be empty or all values are NULL/empty")
	} else {
		fmt.Printf("\n✅ Showed %d results\n", resultCount)
	}

	return nil
}
