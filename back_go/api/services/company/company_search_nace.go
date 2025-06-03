package company

import (
	"context"
	"csv-importer/api/models"
	"fmt"
	"log/slog"
	"time"
)

func (s *companyService) SearchByNaceCode(ctx context.Context, naceCode string, limit int) (*models.CompanySearchResult, error) {
	if naceCode == "" {
		return nil, fmt.Errorf("nace code cannot be empty")
	}

	if limit <= 0 {
		limit = 50
	}

	cacheKey := fmt.Sprintf("companies:full:nace:%s", naceCode)

	start := time.Now()
	var allCompanies []models.CompanyResult
	err := s.cache.Get(cacheKey, &allCompanies)
	cacheDuration := time.Since(start)

	if err != nil {
		slog.Info("Cache miss, fetching complete data from database",
			"nace_code", naceCode,
			"cache_duration_ms", cacheDuration.Milliseconds())

		entityNumbers, err := s.getAllEntityNumbersByNace(naceCode)
		if err != nil {
			return nil, err
		}

		if len(entityNumbers) == 0 {
			return &models.CompanySearchResult{
				Criteria: models.CompanySearchCriteria{NaceCode: naceCode},
				Results:  []models.CompanyResult{},
				Meta:     models.Meta{Count: 0, Total: 0, Limit: limit},
			}, nil
		}

		if len(entityNumbers) > MAX_COMPANIES {
			slog.Warn("Dataset too large, truncating",
				"nace_code", naceCode,
				"original_count", len(entityNumbers),
				"truncated_to", MAX_COMPANIES)
			entityNumbers = entityNumbers[:MAX_COMPANIES]
		}

		allCompanies, err = s.enrichCompleteCompanyData(entityNumbers, naceCode)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Error("Cache write failed", "nace_code", naceCode, "error", err.Error())
		} else {
			slog.Info("Cached complete company dataset", "nace_code", naceCode, "total", len(allCompanies))
		}
	} else {
		slog.Info("Cache hit for complete dataset", "nace_code", naceCode, "total", len(allCompanies))
	}

	total := len(allCompanies)
	var results []models.CompanyResult

	if limit > total {
		results = allCompanies
	} else {
		results = allCompanies[:limit]
	}

	return &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{NaceCode: naceCode},
		Results:  results,
		Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
	}, nil
}

func (s *companyService) getAllEntityNumbersByNace(naceCode string) ([]string, error) {
	query := `
		SELECT DISTINCT entitynumber
		FROM activity
		WHERE nacecode = $1 AND classification = 'MAIN'
		ORDER BY entitynumber
	`

	rows, err := s.db.Query(query, naceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers", "nace_code", naceCode, "count", len(entityNumbers))
	return entityNumbers, nil
}
