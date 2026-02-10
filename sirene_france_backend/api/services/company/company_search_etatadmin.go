package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"time"
)

func (s *companyService) SearchByEtatAdministratif(ctx context.Context, etat string, limit int) (*models.CompanySearchResult, error) {
	cacheKey := fmt.Sprintf("sirene:full:etatadmin:%s", etat)

	var cached models.CompanySearchResult
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		slog.Info("Cache hit", "key", cacheKey)
		return buildSearchResult(&cached, limit), nil
	}

	sirens, err := s.getAllSirensByEtatAdministratif(ctx, etat)
	if err != nil {
		return nil, err
	}

	if len(sirens) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{EtatAdministratif: etat},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	if len(sirens) > MAX_COMPANIES {
		sirens = sirens[:MAX_COMPANIES]
	}

	companies := s.enrichCompanyData(ctx, sirens)

	result := &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{EtatAdministratif: etat},
		Results:  companies,
		Meta:     models.Meta{Total: len(companies)},
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)
	return buildSearchResult(result, limit), nil
}

func (s *companyService) getAllSirensByEtatAdministratif(ctx context.Context, etat string) ([]string, error) {
	query := `SELECT DISTINCT siren FROM unite_legale WHERE etat_administratif_unite_legale = $1`
	rows, err := s.db.QueryContext(ctx, query, etat)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
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
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}
	return sirens, nil
}
