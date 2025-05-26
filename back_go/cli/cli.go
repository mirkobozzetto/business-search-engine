package cli

import (
	"csv-importer/csv"
	"csv-importer/query"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Run(db *sql.DB, args []string) {
	if len(args) < 2 {
		ShowHelp()
		os.Exit(1)
	}

	switch args[1] {
	case "all":
		ProcessAllCSVs(db)
		fmt.Println("🎉 Import done!")
	case "list":
		ListAvailableCSVs()
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
		ShowHelp()
	default:
		ProcessSingleCSV(db, args[1:])
		fmt.Println("🎉 Import done!")
	}
}

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

func ProcessSingleCSV(db *sql.DB, args []string) {
	csvPath := args[0]

	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		tryPath := filepath.Join("../bce_mai_2025", csvPath)
		if _, err := os.Stat(tryPath); err == nil {
			csvPath = tryPath
			fmt.Printf("📁 Found %s in ../bce_mai_2025/\n", args[0])
		} else {
			log.Fatalf("❌ File not found: %s", args[0])
		}
	}

	var tableName string
	if len(args) > 1 {
		tableName = args[1]
	} else {
		tableName = GenerateTableName(csvPath)
	}

	fmt.Printf("🔄 Processing: %s → table '%s'\n", filepath.Base(csvPath), tableName)

	if err := csv.ProcessCSV(db, csvPath, tableName); err != nil {
		log.Fatal("❌ CSV processing failed:", err)
	}
}

func ProcessAllCSVs(db *sql.DB) {
	csvDir := "../bce_mai_2025"
	if err := csv.ProcessAllCSVs(db, csvDir); err != nil {
		log.Fatal("❌ Batch processing failed:", err)
	}
}

func ListAvailableCSVs() {
	csvDir := "../bce_mai_2025"

	files, err := os.ReadDir(csvDir)
	if err != nil {
		log.Fatalf("❌ Cannot read directory %s: %v", csvDir, err)
	}

	fmt.Printf("📁 Available CSV files in %s:\n\n", csvDir)

	csvCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			tableName := GenerateTableName(file.Name())
			info, _ := file.Info()
			size := FormatFileSize(info.Size())
			fmt.Printf("   📄 %-20s → table '%s' (%s)\n", file.Name(), tableName, size)
			csvCount++
		}
	}

	if csvCount == 0 {
		fmt.Println("   No CSV files found")
	}
}

func GenerateTableName(filePath string) string {
	tableName := strings.TrimSuffix(filepath.Base(filePath), ".csv")
	tableName = strings.ReplaceAll(tableName, "-", "_")
	tableName = strings.ReplaceAll(tableName, " ", "_")
	tableName = strings.ToLower(tableName)
	return tableName
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
