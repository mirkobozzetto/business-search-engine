package cli

import (
	"csv-importer/api"
	"csv-importer/cli/_cli"
	"csv-importer/query"
	"database/sql"
	"fmt"
	"os"
)

func Run(db *sql.DB, args []string) {
	if len(args) < 2 {
		_cli.ShowHelp()
		os.Exit(1)
	}

	switch args[1] {
	case "api":
		api.StartAPIServer()
		fmt.Println("🚀 API Server started")
		os.Exit(0)
	case "all":
		_cli.ProcessAllCSVs(db)
		fmt.Println("🎉 Import done!")
	case "list":
		_cli.ListAvailableCSVs()
	case "tables":
		if err := query.ListTables(db); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "stats":
		if err := query.ShowStats(db); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "info":
		if len(args) < 3 {
			fmt.Println("❌ Usage: go run main.go info <table_name>")
			os.Exit(1)
		}
		if err := query.ShowTableInfo(db, args[2]); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "columns":
		if len(args) < 3 {
			fmt.Println("❌ Usage: go run main.go columns <table_name>")
			os.Exit(1)
		}
		if err := query.ShowColumns(db, args[2]); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "values":
		if len(args) < 4 {
			fmt.Println("❌ Usage: go run main.go values <table_name> <column_name> [limit]")
			os.Exit(1)
		}
		limit := 20
		if len(args) > 4 {
			fmt.Sscanf(args[4], "%d", &limit)
		}
		if err := query.ShowColumnValues(db, args[2], args[3], limit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "search":
		if len(args) < 5 {
			fmt.Println("❌ Usage: go run main.go search <table_name> <column_name> <search_value> [limit]")
			os.Exit(1)
		}
		limit := 10
		if len(args) > 5 {
			fmt.Sscanf(args[5], "%d", &limit)
		}
		if err := query.SearchTable(db, args[2], args[3], args[4], limit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "count":
		if len(args) < 5 {
			fmt.Println("❌ Usage: go run main.go count <table_name> <column_name> <search_value>")
			os.Exit(1)
		}
		if err := query.CountRows(db, args[2], args[3], args[4]); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "sample":
		if len(args) < 5 {
			fmt.Println("❌ Usage: go run main.go sample <table_name> <column_name> <search_value> [limit]")
			os.Exit(1)
		}
		limit := 10
		if len(args) > 5 {
			fmt.Sscanf(args[5], "%d", &limit)
		}
		if err := query.SampleRows(db, args[2], args[3], args[4], limit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "export":
		if len(args) < 6 {
			fmt.Println("❌ Usage: go run main.go export <table_name> <column_name> <search_value> <filename.csv>")
			os.Exit(1)
		}
		if err := query.ExportToCSV(db, args[2], args[3], args[4], args[5]); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "preview":
		if len(args) < 3 {
			fmt.Println("❌ Usage: go run main.go preview <table_name> [limit]")
			os.Exit(1)
		}
		limit := 5
		if len(args) > 3 {
			fmt.Sscanf(args[3], "%d", &limit)
		}
		if err := query.PreviewTable(db, args[2], limit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	case "help", "--help", "-h":
		_cli.ShowHelp()
	default:
		_cli.ProcessSingleCSV(db, args[1:])
		fmt.Println("🎉 Import done!")
	}
}
