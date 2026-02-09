package handlers

import (
	"csv-importer/query"
	"database/sql"
	"fmt"
	"os"
)

func HandleListTables(db *sql.DB) {
	if err := query.ListTables(db); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandleShowStats(db *sql.DB) {
	if err := query.ShowStats(db); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandleTableInfo(db *sql.DB, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Usage: go run main.go info <table_name>")
		os.Exit(1)
	}

	tableName := args[0]
	if err := query.ShowTableInfo(db, tableName); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandleShowColumns(db *sql.DB, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Usage: go run main.go columns <table_name>")
		os.Exit(1)
	}

	tableName := args[0]
	if err := query.ShowColumns(db, tableName); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}

func HandleColumnValues(db *sql.DB, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Usage: go run main.go values <table_name> <column_name> [limit]")
		os.Exit(1)
	}

	tableName := args[0]
	columnName := args[1]
	limit := 20

	if len(args) > 2 {
		fmt.Sscanf(args[2], "%d", &limit)
	}

	if err := query.ShowColumnValues(db, tableName, columnName, limit); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
}
