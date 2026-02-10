package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"time"
)

func (s *companyService) SearchMultiCriteria(ctx context.Context, criteria models.CompanySearchCriteria, limit int, offset int) (*models.CompanySearchResult, error) {
	cacheKey := fmt.Sprintf("sirene:multi:%s:%s:%s:%s:%s:%s:%s:%s:%s",
		criteria.NafCode, criteria.Denomination, criteria.CodePostal, criteria.Commune,
		criteria.EtatAdministratif, criteria.DateCreationFrom, criteria.DateCreationTo,
		criteria.CategorieJuridique, criteria.TrancheEffectifs)

	var cached models.CompanySearchResult
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		slog.Info("Cache hit", "key", cacheKey)
		return buildSearchResult(&cached, limit, offset), nil
	}

	conditions := []string{"e.etablissement_siege = 'true'"}
	var args []any
	argN := 1

	if criteria.NafCode != "" {
		conditions = append(conditions, fmt.Sprintf("e.activite_principale_etablissement = $%d", argN))
		args = append(args, criteria.NafCode)
		argN++
	}

	if criteria.Denomination != "" {
		conditions = append(conditions, fmt.Sprintf("u.denomination_unite_legale ILIKE $%d", argN))
		args = append(args, "%"+criteria.Denomination+"%")
		argN++
	}

	if criteria.CodePostal != "" {
		conditions = append(conditions, fmt.Sprintf("e.code_postal_etablissement = $%d", argN))
		args = append(args, criteria.CodePostal)
		argN++
	}

	if criteria.Commune != "" {
		conditions = append(conditions, fmt.Sprintf("e.libelle_commune_etablissement ILIKE $%d", argN))
		args = append(args, "%"+criteria.Commune+"%")
		argN++
	}

	if criteria.EtatAdministratif != "" {
		conditions = append(conditions, fmt.Sprintf("u.etat_administratif_unite_legale = $%d", argN))
		args = append(args, criteria.EtatAdministratif)
		argN++
	}

	if criteria.DateCreationFrom != "" {
		conditions = append(conditions, fmt.Sprintf("u.date_creation_unite_legale >= $%d", argN))
		args = append(args, criteria.DateCreationFrom)
		argN++
	}

	if criteria.DateCreationTo != "" {
		conditions = append(conditions, fmt.Sprintf("u.date_creation_unite_legale <= $%d", argN))
		args = append(args, criteria.DateCreationTo)
		argN++
	}

	if criteria.CategorieJuridique != "" {
		conditions = append(conditions, fmt.Sprintf("u.categorie_juridique_unite_legale = $%d", argN))
		args = append(args, criteria.CategorieJuridique)
		argN++
	}

	if criteria.TrancheEffectifs != "" {
		conditions = append(conditions, fmt.Sprintf("u.tranche_effectifs_unite_legale = $%d", argN))
		args = append(args, criteria.TrancheEffectifs)
		argN++
	}

	if len(conditions) == 1 {
		return &models.CompanySearchResult{
			Criteria: criteria,
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	where := ""
	for i, cond := range conditions {
		if i == 0 {
			where = "WHERE " + cond
		} else {
			where += " AND " + cond
		}
	}

	baseQuery := "SELECT DISTINCT e.siren FROM etablissement e JOIN unite_legale u ON e.siren = u.siren " + where

	countArgs := make([]any, len(args))
	copy(countArgs, args)
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") sub"
	if err := s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	dataQuery := baseQuery + fmt.Sprintf(" LIMIT $%d", argN)
	args = append(args, MAX_COMPANIES)

	rows, err := s.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	var sirens []string
	for rows.Next() {
		var siren string
		if err := rows.Scan(&siren); err != nil {
			continue
		}
		sirens = append(sirens, siren)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(sirens) == 0 {
		return &models.CompanySearchResult{
			Criteria: criteria,
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: totalCount, Count: 0},
		}, nil
	}

	companies := s.enrichCompanyData(ctx, sirens)

	result := &models.CompanySearchResult{
		Criteria: criteria,
		Results:  companies,
		Meta:     models.Meta{Total: totalCount},
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)
	return buildSearchResult(result, limit, offset), nil
}
