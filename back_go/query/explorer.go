package query

import (
	"database/sql"
	"fmt"
	"os"
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

func PreviewTable(db *sql.DB, tableName string, limit int) error {
	if limit <= 0 {
		limit = 5
	}

	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, limit)

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error previewing table: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error getting columns: %v", err)
	}

	fmt.Printf("\n👀 PREVIEW: %s (first %d rows)\n", strings.ToUpper(tableName), limit)

	// Print headers
	for i, col := range columns {
		if i > 0 { fmt.Print(" | ") }
		fmt.Printf("%-15s", col)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", len(columns)*18))

	// Print rows
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		for i, val := range values {
			if i > 0 { fmt.Print(" | ") }

			str := ""
			if val.Valid {
				str = val.String
				if len(str) > 15 {
					str = str[:12] + "..."
				}
			}
			fmt.Printf("%-15s", str)
		}
		fmt.Println()
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

func SearchTable(db *sql.DB, tableName, columnName, searchValue string, limit int) error {
	if limit <= 0 {
		limit = 10
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s ILIKE $1
		LIMIT %d
	`, columnName, tableName, columnName, limit)

	rows, err := db.Query(query, "%"+searchValue+"%")
	if err != nil {
		return fmt.Errorf("error searching table: %v", err)
	}
	defer rows.Close()

	fmt.Printf("\n🔎 SEARCH: %s.%s contains '%s'\n", strings.ToUpper(tableName), strings.ToUpper(columnName), searchValue)
	fmt.Println(strings.Repeat("-", 40))

	count := 0
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			continue
		}
		fmt.Printf("   • %s\n", value)
		count++
	}

	if count == 0 {
		fmt.Println("   No results found.")
	} else {
		fmt.Printf("\n📊 Found %d results\n", count)
	}

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

func CountRows(db *sql.DB, tableName, columnName, searchValue string) error {
	query := fmt.Sprintf(`
		SELECT count(*)
		FROM %s
		WHERE %s ILIKE $1
	`, tableName, columnName)

	var count int64
	err := db.QueryRow(query, "%"+searchValue+"%").Scan(&count)
	if err != nil {
		return fmt.Errorf("error counting rows: %v", err)
	}

	fmt.Printf("\n🔢 COUNT: %s.%s contains '%s'\n", strings.ToUpper(tableName), strings.ToUpper(columnName), searchValue)
	fmt.Printf("📊 Found: %d rows\n", count)

	if count > 1000 {
		fmt.Printf("💡 Tip: Use 'sample' to see examples or 'export' to save to CSV\n")
	}

	return nil
}

func SampleRows(db *sql.DB, tableName, columnName, searchValue string, limit int) error {
	if limit <= 0 {
		limit = 10
	}

	// Get a few columns for context
	columnsQuery := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
		LIMIT 5
	`

	rows, err := db.Query(columnsQuery, tableName)
	if err != nil {
		return fmt.Errorf("error getting columns: %v", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {continue}
		columns = append(columns, col)
	}

	if len(columns) == 0 {
		return fmt.Errorf("no columns found in table %s", tableName)
	}

	// Sample query
	sampleQuery := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s ILIKE $1
		LIMIT %d
	`, strings.Join(columns, ", "), tableName, columnName, limit)

	sampleRows, err := db.Query(sampleQuery, "%"+searchValue+"%")
	if err != nil {
		return fmt.Errorf("error sampling rows: %v", err)
	}
	defer sampleRows.Close()

	fmt.Printf("\n📝 SAMPLE: %s.%s contains '%s' (first %d)\n", strings.ToUpper(tableName), strings.ToUpper(columnName), searchValue, limit)

	// Print headers
	for i, col := range columns {
		if i > 0 { fmt.Print(" | ") }
		fmt.Printf("%-20s", strings.ToUpper(col))
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", len(columns)*23))

	// Print results
	rowCount := 0
	for sampleRows.Next() {
		values := make([]sql.NullString, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := sampleRows.Scan(valuePtrs...); err != nil {
			continue
		}

		for i, val := range values {
			if i > 0 { fmt.Print(" | ") }

			str := ""
			if val.Valid {
				str = val.String
				if len(str) > 20 {
					str = str[:17] + "..."
				}
			}
			fmt.Printf("%-20s", str)
		}
		fmt.Println()
		rowCount++
	}

	if rowCount == 0 {
		fmt.Println("No results found.")
	}

	return nil
}

func ExportToCSV(db *sql.DB, tableName, columnName, searchValue, filename string) error {
	// Get all columns
	columnsQuery := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.Query(columnsQuery, tableName)
	if err != nil {
		return fmt.Errorf("error getting columns: %v", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {continue}
		columns = append(columns, col)
	}

	if len(columns) == 0 {
		return fmt.Errorf("no columns found in table %s", tableName)
	}

	// Count first
	countQuery := fmt.Sprintf(`SELECT count(*) FROM %s WHERE %s ILIKE $1`, tableName, columnName)
	var totalRows int64
	db.QueryRow(countQuery, "%"+searchValue+"%").Scan(&totalRows)

	fmt.Printf("📊 Found %d rows to export...\n", totalRows)

	// Export query
	exportQuery := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s ILIKE $1
	`, strings.Join(columns, ", "), tableName, columnName)

	exportRows, err := db.Query(exportQuery, "%"+searchValue+"%")
	if err != nil {
		return fmt.Errorf("error querying data for export: %v", err)
	}
	defer exportRows.Close()

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write CSV header
	headerLine := strings.Join(columns, ",") + "\n"
	file.WriteString(headerLine)

	// Write data
	rowCount := 0
	for exportRows.Next() {
		values := make([]sql.NullString, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := exportRows.Scan(valuePtrs...); err != nil {
			continue
		}

		var csvValues []string
		for _, val := range values {
			str := ""
			if val.Valid {
				// Escape CSV special characters
				str = strings.ReplaceAll(val.String, "\"", "\"\"")
				if strings.Contains(str, ",") || strings.Contains(str, "\n") || strings.Contains(str, "\"") {
					str = "\"" + str + "\""
				}
			}
			csvValues = append(csvValues, str)
		}

		file.WriteString(strings.Join(csvValues, ",") + "\n")
		rowCount++

		// Progress indicator
		if rowCount%10000 == 0 {
			fmt.Printf("📝 Exported %d rows...\n", rowCount)
		}
	}

	fmt.Printf("✅ EXPORTED: %d rows to %s\n", rowCount, filename)
	fmt.Printf("💡 Open with: Excel, LibreOffice, or any CSV viewer\n")

	return nil
}
