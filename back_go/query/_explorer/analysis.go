package _explorer

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

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
