package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"time"
)

func (s *companyService) SearchByDateCreation(ctx context.Context, fromDate, toDate string, limit int) (*models.CompanySearchResult, error) {
	cacheKey := fmt.Sprintf("sirene:full:datecreation:%s_%s", fromDate, toDate)
	if toDate == "" {
		cacheKey = fmt.Sprintf("sirene:full:datecreation:from_%s", fromDate)
	}

	var cached models.CompanySearchResult
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		slog.Info("Cache hit", "key", cacheKey)
		return buildSearchResult(&cached, limit), nil
	}

	sirens, err := s.getAllSirensByDateCreation(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	if len(sirens) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{DateCreationFrom: fromDate, DateCreationTo: toDate},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	if len(sirens) > MAX_COMPANIES {
		sirens = sirens[:MAX_COMPANIES]
	}

	companies := s.enrichCompanyData(ctx, sirens)

	result := &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{DateCreationFrom: fromDate, DateCreationTo: toDate},
		Results:  companies,
		Meta:     models.Meta{Total: len(companies)},
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)
	return buildSearchResult(result, limit), nil
}

func (s *companyService) getAllSirensByDateCreation(ctx context.Context, fromDate, toDate string) ([]string, error) {
	var query string
	var args []any

	if toDate != "" {
		query = `SELECT DISTINCT siren FROM unite_legale WHERE date_creation_unite_legale >= $1 AND date_creation_unite_legale <= $2`
		args = []any{fromDate, toDate}
	} else {
		query = `SELECT DISTINCT siren FROM unite_legale WHERE date_creation_unite_legale >= $1`
		args = []any{fromDate}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
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
