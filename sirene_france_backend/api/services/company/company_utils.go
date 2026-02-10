package company

import "sirene-importer/api/models"

func buildSearchResult(result *models.CompanySearchResult, limit int) *models.CompanySearchResult {
	if limit <= 0 || limit >= len(result.Results) {
		result.Meta.Count = len(result.Results)
		return result
	}

	return &models.CompanySearchResult{
		Criteria: result.Criteria,
		Results:  result.Results[:limit],
		Meta: models.Meta{
			Total: result.Meta.Total,
			Count: limit,
			Limit: limit,
		},
	}
}
