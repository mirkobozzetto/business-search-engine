package company

import (
	"context"
	"csv-importer/api/models"
	"fmt"
	"log/slog"
	"time"
)

func (s *companyService) SearchByZipcode(ctx context.Context, zipcode string, limit int) (*models.CompanySearchResult, error) {
	if zipcode == "" {
		return nil, fmt.Errorf("zipcode cannot be empty")
	}

	if limit <= 0 {
		limit = 50
	}

	cacheKey := fmt.Sprintf("companies:full:zipcode:%s", zipcode)

	start := time.Now()
	var allCompanies []models.CompanyResult
	err := s.cache.Get(cacheKey, &allCompanies)
	cacheDuration := time.Since(start)

	if err != nil {
		slog.Info("Cache miss, fetching by zipcode from database",
			"zipcode", zipcode,
			"cache_duration_ms", cacheDuration.Milliseconds())

		entityNumbers, err := s.getAllEntityNumbersByZipcode(zipcode)
		if err != nil {
			return nil, err
		}

		if len(entityNumbers) == 0 {
			return &models.CompanySearchResult{
				Criteria: models.CompanySearchCriteria{ZipCode: zipcode},
				Results:  []models.CompanyResult{},
				Meta:     models.Meta{Count: 0, Total: 0, Limit: limit},
			}, nil
		}

		if len(entityNumbers) > MAX_COMPANIES {
			slog.Warn("Dataset too large, truncating",
				"zipcode", zipcode,
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
			slog.Error("Cache write failed", "zipcode", zipcode, "error", err.Error())
		} else {
			slog.Info("Cached complete company dataset", "zipcode", zipcode, "total", len(allCompanies))
		}
	} else {
		slog.Info("Cache hit for complete dataset", "zipcode", zipcode, "total", len(allCompanies))
	}

	return s.buildSearchResult(models.CompanySearchCriteria{ZipCode: zipcode}, allCompanies, limit)
}

func (s *companyService) getAllEntityNumbersByZipcode(zipcode string) ([]string, error) {
	query := `
		SELECT DISTINCT entitynumber
		FROM address
		WHERE zipcode = $1 AND typeofaddress = 'REGO'
		ORDER BY entitynumber
	`

	rows, err := s.db.Query(query, zipcode)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers by zipcode: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers by zipcode", "zipcode", zipcode, "count", len(entityNumbers))
	return entityNumbers, nil
}
