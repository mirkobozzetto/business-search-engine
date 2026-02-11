package csv

import (
	"database/sql"
	"fmt"
	"time"
)

var indexes = []struct {
	name  string
	query string
}{
	{"idx_etab_siren", "CREATE INDEX IF NOT EXISTS idx_etab_siren ON etablissement(siren)"},
	{"idx_etab_siege", "CREATE INDEX IF NOT EXISTS idx_etab_siege ON etablissement(etablissement_siege)"},
	{"idx_etab_naf", "CREATE INDEX IF NOT EXISTS idx_etab_naf ON etablissement(activite_principale_etablissement)"},
	{"idx_etab_cp", "CREATE INDEX IF NOT EXISTS idx_etab_cp ON etablissement(code_postal_etablissement)"},
	{"idx_etab_commune_trgm", "CREATE INDEX IF NOT EXISTS idx_etab_commune_trgm ON etablissement USING gin(libelle_commune_etablissement gin_trgm_ops)"},
	{"idx_etab_siret", "CREATE INDEX IF NOT EXISTS idx_etab_siret ON etablissement(siret)"},
	{"idx_etab_siege_siren", "CREATE INDEX IF NOT EXISTS idx_etab_siege_siren ON etablissement(siren) WHERE etablissement_siege = 'true'"},
	{"idx_etab_siege_naf", "CREATE INDEX IF NOT EXISTS idx_etab_siege_naf ON etablissement(activite_principale_etablissement, siren) WHERE etablissement_siege = 'true'"},
	{"idx_etab_siege_cp", "CREATE INDEX IF NOT EXISTS idx_etab_siege_cp ON etablissement(code_postal_etablissement, siren) WHERE etablissement_siege = 'true'"},
	{"idx_ul_siren", "CREATE INDEX IF NOT EXISTS idx_ul_siren ON unite_legale(siren)"},
	{"idx_ul_etat", "CREATE INDEX IF NOT EXISTS idx_ul_etat ON unite_legale(etat_administratif_unite_legale)"},
	{"idx_ul_date", "CREATE INDEX IF NOT EXISTS idx_ul_date ON unite_legale(date_creation_unite_legale)"},
	{"idx_ul_denom_trgm", "CREATE INDEX IF NOT EXISTS idx_ul_denom_trgm ON unite_legale USING gin(denomination_unite_legale gin_trgm_ops)"},
	{"idx_ul_siren_date", "CREATE INDEX IF NOT EXISTS idx_ul_siren_date ON unite_legale(siren, date_creation_unite_legale DESC)"},
}

func CreateIndexes(db *sql.DB) error {
	fmt.Println("Activation de l'extension pg_trgm...")
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm"); err != nil {
		return fmt.Errorf("pg_trgm: %w", err)
	}

	fmt.Printf("Creation de %d indexes...\n", len(indexes))
	totalStart := time.Now()
	errCount := 0

	for i, idx := range indexes {
		start := time.Now()
		fmt.Printf("[%d/%d] %s... ", i+1, len(indexes), idx.name)

		if _, err := db.Exec(idx.query); err != nil {
			fmt.Printf("ERREUR: %v\n", err)
			errCount++
			continue
		}

		fmt.Printf("OK (%.1fs)\n", time.Since(start).Seconds())
	}

	fmt.Printf("\nIndexes termines en %.1fs (%d/%d reussis)\n", time.Since(totalStart).Seconds(), len(indexes)-errCount, len(indexes))
	if errCount > 0 {
		return fmt.Errorf("%d indexes en echec sur %d", errCount, len(indexes))
	}
	return nil
}
