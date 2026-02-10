package handlers

import (
	"database/sql"
	"fmt"
	"sirene-importer/csv"
)

func HandleCreateIndexes(db *sql.DB) {
	if err := csv.CreateIndexes(db); err != nil {
		fmt.Printf("Erreur: %v\n", err)
	}
}
