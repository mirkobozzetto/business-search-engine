package cli

import (
	"csv-importer/csv"
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
	case "list":
		ListAvailableCSVs()
	case "help", "--help", "-h":
		ShowHelp()
	default:
		ProcessSingleCSV(db, args[1:])
	}

	fmt.Println("🎉 Import done!")
}

func ShowHelp() {
	fmt.Println(`
🚀 CSV Importer - Ultra-flexible CSV to PostgreSQL

Usage:
  go run main.go <csv_path> [table_name]    # Import single CSV
  go run main.go all                        # Import all CSVs from ../bce_mai_2025/
  go run main.go list                       # List available CSV files
  go run main.go help                       # Show this help

Examples:
  go run main.go ../data/users.csv                    # users.csv → table 'users'
  go run main.go ../data/users.csv custom_users       # users.csv → table 'custom_users'
  go run main.go activity.csv                         # looks in ../bce_mai_2025/activity.csv
  go run main.go all                                  # Import all CSV files
  go run main.go list                                 # See what CSV files are available

Features:
  ⚡ 1.3M+ lines/sec using PostgreSQL COPY
  🔥 Automatic table name generation
  📊 Real-time progress monitoring
  🏗️ Clean table structure from CSV headers`)
}

func ProcessSingleCSV(db *sql.DB, args []string) {
	csvPath := args[0]

	// If path doesn't exist, try in bce_mai_2025 directory
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

			// Get file size
			info, _ := file.Info()
			size := FormatFileSize(info.Size())

			fmt.Printf("   📄 %-20s → table '%s' (%s)\n", file.Name(), tableName, size)
			csvCount++
		}
	}

	if csvCount == 0 {
		fmt.Println("   No CSV files found")
	} else {
		fmt.Printf("\n💡 Use 'go run main.go all' to import all %d files\n", csvCount)
		fmt.Printf("💡 Use 'go run main.go <filename>' to import individual files\n")
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
