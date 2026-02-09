package company

import (
	"context"
	"csv-importer/api/models"
	"fmt"
	"log/slog"
)

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

	if criteria.StartDateFrom != "" {
		var cacheKey string
		if criteria.StartDateTo != "" {
			cacheKey = fmt.Sprintf("companies:full:startdate:%s_%s", criteria.StartDateFrom, criteria.StartDateTo)
		} else {
			cacheKey = fmt.Sprintf("companies:full:startdate:from_%s", criteria.StartDateFrom)
		}

		var companies []models.CompanyResult
		err := s.cache.Get(cacheKey, &companies)
		if err != nil {
			return nil, fmt.Errorf("start date cache not found: %s. Please search by start date first", cacheKey)
		}
		allDatasets = append(allDatasets, companies)
		criteriaCount++
		slog.Info("Found cached start date data", "from", criteria.StartDateFrom, "to", criteria.StartDateTo, "count", len(companies))
	}

	if criteriaCount == 0 {
		return nil, fmt.Errorf("at least one search criteria required")
	}

	if criteriaCount == 1 {
		return s.buildSearchResult(criteria, allDatasets[0], limit)
	}

	intersection := s.intersectCompanyResults(allDatasets)

	slog.Info("Multi-criteria intersection",
		"criteria_count", criteriaCount,
		"datasets_sizes", fmt.Sprintf("%v", getDatasetSizes(allDatasets)),
		"intersection_size", len(intersection))

	return s.buildSearchResult(criteria, intersection, limit)
}
