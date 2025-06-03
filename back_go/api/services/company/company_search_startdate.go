package company

import (
	"context"
	"csv-importer/api/models"
	"fmt"
	"log/slog"
	"time"
)

func (s *companyService) SearchByStartDate(ctx context.Context, fromDate, toDate string, limit int) (*models.CompanySearchResult, error) {
	if fromDate == "" {
		return nil, fmt.Errorf("start date from cannot be empty")
	}

	if limit <= 0 {
		limit = 50
	}

	var cacheKey string
	if toDate != "" {
		cacheKey = fmt.Sprintf("companies:full:startdate:%s_%s", fromDate, toDate)
	} else {
		cacheKey = fmt.Sprintf("companies:full:startdate:from_%s", fromDate)
	}

	start := time.Now()
	var allCompanies []models.CompanyResult
	err := s.cache.Get(cacheKey, &allCompanies)
	cacheDuration := time.Since(start)

	if err != nil {
		slog.Info("Cache miss, fetching by start date from database",
			"from_date", fromDate,
			"to_date", toDate,
			"cache_duration_ms", cacheDuration.Milliseconds())

		entityNumbers, err := s.getAllEntityNumbersByStartDate(fromDate, toDate)
		if err != nil {
			return nil, err
		}

		if len(entityNumbers) == 0 {
			return &models.CompanySearchResult{
				Criteria: models.CompanySearchCriteria{StartDateFrom: fromDate, StartDateTo: toDate},
				Results:  []models.CompanyResult{},
				Meta:     models.Meta{Count: 0, Total: 0, Limit: limit},
			}, nil
		}

		if len(entityNumbers) > MAX_COMPANIES {
			slog.Warn("Dataset too large, truncating",
				"from_date", fromDate,
				"to_date", toDate,
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
			slog.Error("Cache write failed", "from_date", fromDate, "to_date", toDate, "error", err.Error())
		} else {
			slog.Info("Cached complete company dataset", "from_date", fromDate, "to_date", toDate, "total", len(allCompanies))
		}
	} else {
		slog.Info("Cache hit for complete dataset", "from_date", fromDate, "to_date", toDate, "total", len(allCompanies))
	}

	return s.buildSearchResult(models.CompanySearchCriteria{StartDateFrom: fromDate, StartDateTo: toDate}, allCompanies, limit)
}

func (s *companyService) getAllEntityNumbersByStartDate(fromDate, toDate string) ([]string, error) {
	var query string
	var args []any

	if fromDate != "" && toDate != "" {
		query = `
			SELECT DISTINCT enterprisenumber
			FROM enterprise
			WHERE TO_DATE(startdate, 'DD-MM-YYYY') >= TO_DATE($1, 'DD-MM-YYYY')
			  AND TO_DATE(startdate, 'DD-MM-YYYY') <= TO_DATE($2, 'DD-MM-YYYY')
			ORDER BY enterprisenumber
		`
		args = []any{fromDate, toDate}
		slog.Info("Searching by date range", "from", fromDate, "to", toDate)
	} else if fromDate != "" {
		query = `
			SELECT DISTINCT enterprisenumber
			FROM enterprise
			WHERE TO_DATE(startdate, 'DD-MM-YYYY') >= TO_DATE($1, 'DD-MM-YYYY')
			ORDER BY enterprisenumber
		`
		args = []any{fromDate}
		slog.Info("Searching by start date", "from", fromDate)
	} else {
		return nil, fmt.Errorf("at least fromDate is required")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers by start date: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers by start date", "from", fromDate, "to", toDate, "count", len(entityNumbers))
	return entityNumbers, nil
}
