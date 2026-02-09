package handlers

import "fmt"

func ShowHelp() {
	fmt.Print(`
🔧 CSV IMPORTER - BCE Database Tool

USAGE:
  go run main.go <command> [arguments]

COMMANDS:
  📊 DATABASE OPERATIONS:
    api                              Launch API server
    all                             Import all CSV files in parallel
    list                            List available CSV files

  📋 TABLE MANAGEMENT:
    tables                          List all database tables
    stats                           Show database statistics
    info <table_name>               Show table information
    columns <table_name>            Show table columns
    preview <table_name> [limit]    Preview table data

  🔍 SEARCH & ANALYSIS:
    search <table> <column> <value> [limit]    Search in table
    count <table> <column> <value>             Count matching rows
    sample <table> <column> <value> [limit]    Sample matching rows
    values <table> <column> [limit]            Show unique values

  📤 EXPORT:
    export <table> <column> <value> <file.csv> Export search results

  ℹ️  HELP:
    help, --help, -h                Show this help

EXAMPLES:
  go run main.go api                                    # Start API server
  go run main.go all                                    # Import all CSVs
  go run main.go search activity nacecode 62020 100    # Find companies in IT sector
  go run main.go export activity nacecode 62020 it.csv # Export IT companies
  go run main.go preview denomination 10                # Preview company names

🚀 For more info: https://github.com/yourusername/csv-importer
`)
}
