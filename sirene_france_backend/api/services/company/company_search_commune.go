package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"time"
)

func (s *companyService) SearchByCommune(ctx context.Context, commune string, limit, offset int) (*models.CompanySearchResult, error) {
	cacheKey := fmt.Sprintf("sirene:full:commune:%s", commune)

	var cached models.CompanySearchResult
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		slog.Info("Cache hit", "key", cacheKey)
		return buildSearchResult(&cached, limit, offset), nil
	}

	sirens, err := s.getAllSirensByCommune(ctx, commune)
	if err != nil {
		return nil, err
	}

	if len(sirens) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{Commune: commune},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	totalFound := len(sirens)
	if len(sirens) > MAX_COMPANIES {
		sirens = sirens[:MAX_COMPANIES]
	}

	companies := s.enrichCompanyData(ctx, sirens)

	result := &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{Commune: commune},
		Results:  companies,
		Meta:     models.Meta{Total: totalFound},
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)
	return buildSearchResult(result, limit, offset), nil
}

func (s *companyService) getAllSirensByCommune(ctx context.Context, commune string) ([]string, error) {
	query := `SELECT DISTINCT siren FROM etablissement WHERE libelle_commune_etablissement ILIKE $1 AND etablissement_siege = 'true'`
	rows, err := s.db.QueryContext(ctx, query, "%"+commune+"%")
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
