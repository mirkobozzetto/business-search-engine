package handlers

import (
	"csv-importer/csv"
	"database/sql"
	"log"
)

func HandleImportAll(db *sql.DB) {
	csvDir := "../bce_mai_2025"
	if err := csv.ProcessAllCSVsParallel(db, csvDir); err != nil {
		log.Fatal("❌ Parallel batch processing failed:", err)
	}
}

func HandleListCSVs() {
	// TODO: Move logic from _cli/list.go here
	csvDir := "../bce_mai_2025"
	// Implementation will be moved here
	log.Printf("📁 Listing CSV files in %s", csvDir)
}
