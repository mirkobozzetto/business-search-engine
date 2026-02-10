package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"strings"
)

func (s *companyService) enrichCompanyData(ctx context.Context, sirens []string) []models.CompanyResult {
	companyMap := make(map[string]*models.CompanyResult)
	for _, siren := range sirens {
		companyMap[siren] = &models.CompanyResult{Siren: siren}
	}

	s.enrichFromUniteLegale(ctx, companyMap, sirens)
	s.enrichFromEtablissement(ctx, companyMap, sirens)

	results := make([]models.CompanyResult, 0, len(companyMap))
	for _, company := range companyMap {
		results = append(results, *company)
	}
	return results
}

func (s *companyService) enrichFromUniteLegale(ctx context.Context, companyMap map[string]*models.CompanyResult, sirens []string) {
	batchSize := 1000
	for i := 0; i < len(sirens); i += batchSize {
		end := i + batchSize
		if end > len(sirens) {
			end = len(sirens)
		}
		batch := sirens[i:end]

		placeholders := make([]string, len(batch))
		args := make([]any, len(batch))
		for j, siren := range batch {
			placeholders[j] = fmt.Sprintf("$%d", j+1)
			args[j] = siren
		}

		query := fmt.Sprintf(`SELECT siren, denomination_unite_legale, sigle_unite_legale,
			categorie_juridique_unite_legale, date_creation_unite_legale,
			etat_administratif_unite_legale, tranche_effectifs_unite_legale,
			categorie_entreprise, activite_principale_unite_legale
			FROM unite_legale WHERE siren IN (%s)`, strings.Join(placeholders, ","))

		rows, err := s.db.QueryContext(ctx, query, args...)
		if err != nil {
			slog.Error("Enrich unite_legale failed", "error", err)
			continue
		}

		for rows.Next() {
			var siren, denom, sigle, catJuridique, dateCreation, etat, effectifs, catEntreprise, naf string
			err := rows.Scan(&siren, &denom, &sigle, &catJuridique, &dateCreation, &etat, &effectifs, &catEntreprise, &naf)
			if err != nil {
				slog.Warn("Scan unite_legale failed", "error", err)
				continue
			}
			if c, ok := companyMap[siren]; ok {
				c.Denomination = denom
				c.Sigle = sigle
				c.CategorieJuridique = catJuridique
				c.DateCreation = dateCreation
				c.EtatAdministratif = etat
				c.TrancheEffectifs = effectifs
				c.CategorieEntreprise = catEntreprise
				c.NafCode = naf
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			slog.Error("rows iteration error unite_legale", "error", err)
		}
	}
}

func (s *companyService) enrichFromEtablissement(ctx context.Context, companyMap map[string]*models.CompanyResult, sirens []string) {
	batchSize := 1000
	for i := 0; i < len(sirens); i += batchSize {
		end := i + batchSize
		if end > len(sirens) {
			end = len(sirens)
		}
		batch := sirens[i:end]

		placeholders := make([]string, len(batch))
		args := make([]any, len(batch))
		for j, siren := range batch {
			placeholders[j] = fmt.Sprintf("$%d", j+1)
			args[j] = siren
		}

		query := fmt.Sprintf(`SELECT siren, siret, enseigne1_etablissement,
			numero_voie_etablissement, type_voie_etablissement, libelle_voie_etablissement,
			code_postal_etablissement, libelle_commune_etablissement, activite_principale_etablissement
			FROM etablissement WHERE siren IN (%s) AND etablissement_siege = 'true'`, strings.Join(placeholders, ","))

		rows, err := s.db.QueryContext(ctx, query, args...)
		if err != nil {
			slog.Error("Enrich etablissement failed", "error", err)
			continue
		}

		for rows.Next() {
			var siren, siret, enseigne, numVoie, typeVoie, libelleVoie, cp, commune, nafEtab string
			err := rows.Scan(&siren, &siret, &enseigne, &numVoie, &typeVoie, &libelleVoie, &cp, &commune, &nafEtab)
			if err != nil {
				slog.Warn("Scan etablissement failed", "error", err)
				continue
			}
			if c, ok := companyMap[siren]; ok {
				c.Siret = siret
				c.Enseigne = enseigne
				c.NumeroVoie = numVoie
				c.TypeVoie = typeVoie
				c.LibelleVoie = libelleVoie
				c.CodePostal = cp
				c.LibelleCommune = commune
				if nafEtab != "" {
					c.NafCode = nafEtab
				}
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			slog.Error("rows iteration error etablissement", "error", err)
		}
	}
}
