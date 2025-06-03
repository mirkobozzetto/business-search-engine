package company

import (
	"context"
	"csv-importer/api/models"
	"fmt"
	"log/slog"
	"time"
)

func (s *companyService) SearchByDenomination(ctx context.Context, query string, limit int) (*models.CompanySearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("denomination query cannot be empty")
	}

	if limit <= 0 {
		limit = 50
	}

	cacheKey := fmt.Sprintf("companies:full:denomination:%s", query)

	start := time.Now()
	var allCompanies []models.CompanyResult
	err := s.cache.Get(cacheKey, &allCompanies)
	cacheDuration := time.Since(start)

	if err != nil {
		slog.Info("Cache miss, fetching by denomination from database",
			"query", query,
			"cache_duration_ms", cacheDuration.Milliseconds())

		entityNumbers, err := s.getAllEntityNumbersByDenomination(query)
		if err != nil {
			return nil, err
		}

		if len(entityNumbers) == 0 {
			return &models.CompanySearchResult{
				Criteria: models.CompanySearchCriteria{Denomination: query},
				Results:  []models.CompanyResult{},
				Meta:     models.Meta{Count: 0, Total: 0, Limit: limit},
			}, nil
		}

		if len(entityNumbers) > MAX_COMPANIES {
			slog.Warn("Dataset too large, truncating",
				"query", query,
				"original_count", len(entityNumbers),
				"truncated_to", MAX_COMPANIES)
			entityNumbers = entityNumbers[:MAX_COMPANIES]
		}

		allCompanies, err = s.enrichCompleteCompanyData(entityNumbers, "")
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Error("Cache write failed", "query", query, "error", err.Error())
		} else {
			slog.Info("Cached complete company dataset", "query", query, "total", len(allCompanies))
		}
	} else {
		slog.Info("Cache hit for complete dataset", "query", query, "total", len(allCompanies))
	}

	return s.buildSearchResult(models.CompanySearchCriteria{Denomination: query}, allCompanies, limit)
}

func (s *companyService) getAllEntityNumbersByDenomination(query string) ([]string, error) {
	querySQL := `
		SELECT DISTINCT entitynumber
		FROM denomination
		WHERE denomination ILIKE $1
		ORDER BY entitynumber
	`

	rows, err := s.db.Query(querySQL, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers by denomination: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers by denomination", "query", query, "count", len(entityNumbers))
	return entityNumbers, nil
}
