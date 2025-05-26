package _cli

import (
	"csv-importer/csv"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

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
