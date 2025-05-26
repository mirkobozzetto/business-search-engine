package _cli

import "fmt"

func ShowHelp() {
	fmt.Println(`
🚀 CSV Importer - Simple & Fast

IMPORT:
  go run main.go <csv_path> [table_name]    # Import single CSV
  go run main.go all                        # Import all CSVs
  go run main.go list                       # List available CSV files

EXPLORE:
  go run main.go tables                     # Show created tables
  go run main.go info <table>               # Table structure
  go run main.go preview <table> [limit]    # Preview data
  go run main.go values <table> <column>    # Column values

ANALYZE:
  go run main.go count <table> <column> <term>      # Count matching rows
  go run main.go sample <table> <column> <term>     # Sample matching rows
  go run main.go export <table> <column> <term> <file.csv>  # Export to CSV

Examples:
  go run main.go all                        # Import everything
  go run main.go values activity nace_code  # See activity codes
  go run main.go count activity nace_code "62020"   # Count IT companies
  go run main.go sample activity nace_code "62020" 5  # See 5 examples
  go run main.go export activity nace_code "62020" it_companies.csv  # Export all`)
}
