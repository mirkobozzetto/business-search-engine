package handlers

import (
	"csv-importer/query"
	"database/sql"
	"fmt"
	"os"
)

func HandleExport(db *sql.DB, args []string) {
	if len(args) < 4 {
		fmt.Println("❌ Usage: go run main.go export <table_name> <column_name> <search_value> <filename.csv>")
		os.Exit(1)
	}

	tableName := args[0]
	columnName := args[1]
	searchValue := args[2]
	filename := args[3]

	if err := query.ExportToCSV(db, tableName, columnName, searchValue, filename); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandlePreview(db *sql.DB, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Usage: go run main.go preview <table_name> [limit]")
		os.Exit(1)
	}

	tableName := args[0]
	limit := 5

	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &limit)
	}

	if err := query.PreviewTable(db, tableName, limit); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}
