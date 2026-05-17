package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

const MAX_COMPANIES = 100000

func (s *companyService) SearchByNafCode(ctx context.Context, nafCode string, limit, offset int) (*models.CompanySearchResult, error) {
	conditions := []string{
		"e.etablissement_siege = 'true'",
		"e.activite_principale_etablissement = $1",
	}
	args := []any{nafCode}
	cacheKey := fmt.Sprintf("sirene:v2:naf:%s", nafCode)
	criteria := models.CompanySearchCriteria{NafCode: nafCode}

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
