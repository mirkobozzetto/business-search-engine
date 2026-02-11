package csv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

type nafSection struct {
	Code  string    `json:"code"`
	Label string    `json:"label"`
	Codes []nafCode `json:"codes"`
}

type nafCode struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type nafFile struct {
	Sections []nafSection `json:"sections"`
}

func LoadNafCodes(db *sql.DB, jsonPath string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS naf_reference (
		code TEXT PRIMARY KEY,
		label TEXT NOT NULL,
		section_code TEXT NOT NULL,
		section_label TEXT NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	_, err = db.Exec("TRUNCATE TABLE naf_reference")
	if err != nil {
		return fmt.Errorf("truncate: %w", err)
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var naf nafFile
	if err := json.Unmarshal(data, &naf); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	count := 0
	for _, section := range naf.Sections {
		for _, code := range section.Codes {
			_, err := db.Exec(
				"INSERT INTO naf_reference (code, label, section_code, section_label) VALUES ($1, $2, $3, $4)",
				code.Code, code.Label, section.Code, section.Label,
			)
			if err != nil {
				return fmt.Errorf("insert %s: %w", code.Code, err)
			}
			count++
		}
	}

	fmt.Printf("%d codes NAF inseres\n", count)
	return nil
}
