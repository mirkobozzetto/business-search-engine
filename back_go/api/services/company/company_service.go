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
		Criteria: models.CompanySearchCriteria{Denomination: query},
		Results:  results,
		Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
	}, nil
}

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

		dataSize := len(allCompanies)
		estimatedSizeMB := dataSize * 2000 / 1024 / 1024

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Error("Cache write failed",
				"zipcode", zipcode,
				"companies_count", dataSize,
				"estimated_size_mb", estimatedSizeMB,
				"error", err.Error())
		} else {
			slog.Info("Cached complete company dataset",
				"zipcode", zipcode,
				"total", len(allCompanies),
				"estimated_size_mb", estimatedSizeMB)
		}
	} else {
		slog.Info("Cache hit for complete dataset",
			"zipcode", zipcode,
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
		Criteria: models.CompanySearchCriteria{ZipCode: zipcode},
		Results:  results,
		Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
	}, nil
}

func (s *companyService) SearchMultiCriteria(ctx context.Context, criteria models.CompanySearchCriteria, limit int) (*models.CompanySearchResult, error) {
	if limit <= 0 {
		limit = 50
	}

	var allDatasets [][]models.CompanyResult
	var criteriaCount int

	if criteria.NaceCode != "" {
		cacheKey := fmt.Sprintf("companies:full:nace:%s", criteria.NaceCode)
		var companies []models.CompanyResult
		err := s.cache.Get(cacheKey, &companies)
		if err != nil {
			return nil, fmt.Errorf("NACE cache not found: %s. Please search by NACE first", criteria.NaceCode)
		}
		allDatasets = append(allDatasets, companies)
		criteriaCount++
		slog.Info("Found cached NACE data", "nace_code", criteria.NaceCode, "count", len(companies))
	}

	if criteria.Denomination != "" {
		cacheKey := fmt.Sprintf("companies:full:denomination:%s", criteria.Denomination)
		var companies []models.CompanyResult
		err := s.cache.Get(cacheKey, &companies)
		if err != nil {
			return nil, fmt.Errorf("denomination cache not found: %s. Please search by denomination first", criteria.Denomination)
		}
		allDatasets = append(allDatasets, companies)
		criteriaCount++
		slog.Info("Found cached denomination data", "denomination", criteria.Denomination, "count", len(companies))
	}

	if criteria.ZipCode != "" {
		cacheKey := fmt.Sprintf("companies:full:zipcode:%s", criteria.ZipCode)
		var companies []models.CompanyResult
		err := s.cache.Get(cacheKey, &companies)
		if err != nil {
			return nil, fmt.Errorf("zipcode cache not found: %s. Please search by zipcode first", criteria.ZipCode)
		}
		allDatasets = append(allDatasets, companies)
		criteriaCount++
		slog.Info("Found cached zipcode data", "zipcode", criteria.ZipCode, "count", len(companies))
	}

	if criteriaCount == 0 {
		return nil, fmt.Errorf("at least one search criteria required")
	}

	if criteriaCount == 1 {
		results := allDatasets[0]
		total := len(results)

		if limit > total {
			results = results[:total]
		} else {
			results = results[:limit]
		}

		return &models.CompanySearchResult{
			Criteria: criteria,
			Results:  results,
			Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
		}, nil
	}

	intersection := s.intersectCompanyResults(allDatasets)

	slog.Info("Multi-criteria intersection",
		"criteria_count", criteriaCount,
		"datasets_sizes", fmt.Sprintf("%v", getDatasetSizes(allDatasets)),
		"intersection_size", len(intersection))

	total := len(intersection)
	var results []models.CompanyResult

	if limit > total {
		results = intersection
	} else {
		results = intersection[:limit]
	}

	return &models.CompanySearchResult{
		Criteria: criteria,
		Results:  results,
		Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
	}, nil
}

func (s *companyService) intersectCompanyResults(datasets [][]models.CompanyResult) []models.CompanyResult {
	if len(datasets) == 0 {
		return []models.CompanyResult{}
	}

	if len(datasets) == 1 {
		return datasets[0]
	}

	smallest := 0
	for i, dataset := range datasets {
		if len(dataset) < len(datasets[smallest]) {
			smallest = i
		}
	}

	entityMap := make(map[string]models.CompanyResult)
	for _, company := range datasets[smallest] {
		entityMap[company.EntityNumber] = company
	}

	for i, dataset := range datasets {
		if i == smallest {
			continue
		}

		newEntityMap := make(map[string]models.CompanyResult)
		for _, company := range dataset {
			if _, exists := entityMap[company.EntityNumber]; exists {
				newEntityMap[company.EntityNumber] = company
			}
		}
		entityMap = newEntityMap
	}

	var result []models.CompanyResult
	for _, company := range entityMap {
		result = append(result, company)
	}

	return result
}

func getDatasetSizes(datasets [][]models.CompanyResult) []int {
	sizes := make([]int, len(datasets))
	for i, dataset := range datasets {
		sizes[i] = len(dataset)
	}
	return sizes
}
