package company

import "sirene-importer/api/models"

func buildSearchResult(result *models.CompanySearchResult, limit, offset int) *models.CompanySearchResult {
	total := len(result.Results)

	if offset >= total {
		pages := 1
		if limit > 0 {
			pages = (total + limit - 1) / limit
		}
		return &models.CompanySearchResult{
			Criteria: result.Criteria,
			Results:  []models.CompanyResult{},
			Meta: models.Meta{
				Total:  total,
				Count:  0,
				Limit:  limit,
				Offset: offset,
				Page:   pages + 1,
				Pages:  pages,
			},
		}
	}

	end := offset + limit
	if limit <= 0 || end > total {
		end = total
	}

	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}
	pages := 1
	if limit > 0 {
		pages = (total + limit - 1) / limit
	}

	return &models.CompanySearchResult{
		Criteria: result.Criteria,
		Results:  result.Results[offset:end],
		Meta: models.Meta{
			Total:  total,
			Count:  end - offset,
			Limit:  limit,
			Offset: offset,
			Page:   page,
			Pages:  pages,
		},
	}
}
