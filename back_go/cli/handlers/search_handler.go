package handlers

import (
	"csv-importer/query"
	"database/sql"
	"fmt"
	"os"
)

func HandleSearch(db *sql.DB, args []string) {
	if len(args) < 3 {
		fmt.Println("❌ Usage: go run main.go search <table_name> <column_name> <search_value> [limit]")
		os.Exit(1)
	}

	tableName := args[0]
	columnName := args[1]
	searchValue := args[2]
	limit := 10

	if len(args) > 3 {
		fmt.Sscanf(args[3], "%d", &limit)
	}

	if err := query.SearchTable(db, tableName, columnName, searchValue, limit); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandleCount(db *sql.DB, args []string) {
	if len(args) < 3 {
		fmt.Println("❌ Usage: go run main.go count <table_name> <column_name> <search_value>")
		os.Exit(1)
	}

	tableName := args[0]
	columnName := args[1]
	searchValue := args[2]

	if err := query.CountRows(db, tableName, columnName, searchValue); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandleSample(db *sql.DB, args []string) {
	if len(args) < 3 {
		fmt.Println("❌ Usage: go run main.go sample <table_name> <column_name> <search_value> [limit]")
		os.Exit(1)
	}

	tableName := args[0]
	columnName := args[1]
	searchValue := args[2]
	limit := 10

	if len(args) > 3 {
		fmt.Sscanf(args[3], "%d", &limit)
	}

	if err := query.SampleRows(db, tableName, columnName, searchValue, limit); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}
