package _explorer

import (
	"database/sql"
	"fmt"
	"strings"
)

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
