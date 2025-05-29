package handlers

import (
	"csv-importer/csv"
	"database/sql"
	"log/slog"
)

func HandleImportAll(db *sql.DB) {
	csvDir := "../bce_mai_2025"
	if err := csv.ProcessAllCSVsParallel(db, csvDir); err != nil {
		slog.Error("❌ Parallel batch processing failed", "error", err)
	}
}

func HandleListCSVs() {
	// TODO: Move logic from _cli/list.go here
	csvDir := "../bce_mai_2025"
	// Implementation will be moved here
	slog.Info("📁 Listing CSV files in", "directory", csvDir)
}
