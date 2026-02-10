package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"time"
)

const MAX_COMPANIES = 100000

func (s *companyService) SearchByNafCode(ctx context.Context, nafCode string, limit, offset int) (*models.CompanySearchResult, error) {
	cacheKey := fmt.Sprintf("sirene:full:naf:%s", nafCode)

	var cached models.CompanySearchResult
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		slog.Info("Cache hit", "key", cacheKey, "total", cached.Meta.Total)
		return buildSearchResult(&cached, limit, offset), nil
	}

	sirens, err := s.getAllSirensByNaf(ctx, nafCode)
	if err != nil {
		return nil, err
	}

	if len(sirens) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{NafCode: nafCode},
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
		Criteria: models.CompanySearchCriteria{NafCode: nafCode},
		Results:  companies,
		Meta:     models.Meta{Total: totalFound},
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)
	return buildSearchResult(result, limit, offset), nil
}

func (s *companyService) getAllSirensByNaf(ctx context.Context, nafCode string) ([]string, error) {
	query := `SELECT DISTINCT siren FROM etablissement WHERE activite_principale_etablissement = $1`
	rows, err := s.db.QueryContext(ctx, query, nafCode)
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
