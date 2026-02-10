package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"strings"
	"time"
)

func (s *companyService) SearchByDenomination(ctx context.Context, query string, limit, offset int) (*models.CompanySearchResult, error) {
	cacheKey := fmt.Sprintf("sirene:full:denomination:%s", query)

	var cached models.CompanySearchResult
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		slog.Info("Cache hit", "key", cacheKey)
		return buildSearchResult(&cached, limit, offset), nil
	}

	sirens, err := s.getAllSirensByDenomination(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(sirens) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{Denomination: query},
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
		Criteria: models.CompanySearchCriteria{Denomination: query},
		Results:  companies,
		Meta:     models.Meta{Total: totalFound},
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)
	return buildSearchResult(result, limit, offset), nil
}

func (s *companyService) getAllSirensByDenomination(ctx context.Context, query string) ([]string, error) {
	words := strings.Fields(strings.TrimSpace(query))
	if len(words) == 0 {
		return nil, fmt.Errorf("empty query")
	}

	var conditions []string
	var args []any
	for i, word := range words {
		args = append(args, "%"+word+"%")
		conditions = append(conditions, fmt.Sprintf("denomination_unite_legale ILIKE $%d", i+1))
	}

	q := fmt.Sprintf(`SELECT DISTINCT siren FROM unite_legale WHERE %s`, strings.Join(conditions, " AND "))
	rows, err := s.db.QueryContext(ctx, q, args...)
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
