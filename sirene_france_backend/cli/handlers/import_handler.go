package handlers

import (
	"database/sql"
	"fmt"
	"sirene-importer/csv"
)

func HandleImportAll(db *sql.DB) {
	fmt.Println("Importing SIRENE ZIP files...")
	if err := csv.ProcessAllZIPs(db, "../sirene_data"); err != nil {
		fmt.Printf("Import error: %v\n", err)
	}
}
