package handlers

import (
	"database/sql"
	"fmt"
	"sirene-importer/csv"
)

func HandleImportNaf(db *sql.DB) {
	fmt.Println("Importing NAF codes from data/naf_codes.json...")
	if err := csv.LoadNafCodes(db, "data/naf_codes.json"); err != nil {
		fmt.Printf("NAF import error: %v\n", err)
	}
}
