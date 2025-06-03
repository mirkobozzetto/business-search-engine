package company

import (
	"csv-importer/api/models"
)

func (s *companyService) buildSearchResult(criteria models.CompanySearchCriteria, allCompanies []models.CompanyResult, limit int) (*models.CompanySearchResult, error) {
	total := len(allCompanies)
	var results []models.CompanyResult

	if limit > total {
		results = allCompanies
	} else {
		results = allCompanies[:limit]
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
