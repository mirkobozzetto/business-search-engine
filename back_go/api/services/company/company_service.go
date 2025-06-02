package company

import (
	"context"
	"csv-importer/api/cache"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const MAX_COMPANIES = 100000

type companyService struct {
	db    *sql.DB
	cache *cache.RedisCache
}

func NewCompanyService(db *sql.DB) CompanyService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	// ! \\ TODO: move to env and add to config
	redisCache := cache.NewRedisCache(cache.CacheConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	})

	return &companyService{
		db:    db,
		cache: redisCache,
	}
}

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

		dataSize := len(allCompanies)
		estimatedSizeMB := dataSize * 2000 / 1024 / 1024

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Error("Cache write failed",
				"nace_code", naceCode,
				"companies_count", dataSize,
				"estimated_size_mb", estimatedSizeMB,
				"error", err.Error())
		} else {
			slog.Info("Cached complete company dataset",
				"nace_code", naceCode,
				"total", len(allCompanies),
				"estimated_size_mb", estimatedSizeMB)
		}
	} else {
		slog.Info("Cache hit for complete dataset",
			"nace_code", naceCode,
			"total", len(allCompanies),
			"cache_duration_ms", cacheDuration.Milliseconds())
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
				Criteria: models.CompanySearchCriteria{},
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

		dataSize := len(allCompanies)
		estimatedSizeMB := dataSize * 2000 / 1024 / 1024

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Error("Cache write failed",
				"query", query,
				"companies_count", dataSize,
				"estimated_size_mb", estimatedSizeMB,
				"error", err.Error())
		} else {
			slog.Info("Cached complete company dataset",
				"query", query,
				"total", len(allCompanies),
				"estimated_size_mb", estimatedSizeMB)
		}
	} else {
		slog.Info("Cache hit for complete dataset",
			"query", query,
			"total", len(allCompanies),
			"cache_duration_ms", cacheDuration.Milliseconds())
	}

	total := len(allCompanies)
	var results []models.CompanyResult

	if limit > total {
		results = allCompanies
	} else {
		results = allCompanies[:limit]
	}

	return &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{},
		Results:  results,
		Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
	}, nil
}
