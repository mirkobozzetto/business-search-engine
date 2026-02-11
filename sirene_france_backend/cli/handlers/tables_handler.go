package handlers

import (
	"database/sql"
	"fmt"
)

func HandleListTables(db *sql.DB) {
	rows, err := db.Query(`
		SELECT table_name,
			   (SELECT COUNT(*) FROM information_schema.columns WHERE columns.table_name = tables.table_name) as col_count
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name`)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer func() { _ = rows.Close() }()

	fmt.Printf("%-30s %s\n", "TABLE", "COLUMNS")
	fmt.Println("-------------------------------------------")

	for rows.Next() {
		var name string
		var cols int
		_ = rows.Scan(&name, &cols)
		fmt.Printf("%-30s %d\n", name, cols)
	}
}
