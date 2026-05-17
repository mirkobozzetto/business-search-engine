package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

func (s *companyService) SearchByCodePostal(ctx context.Context, codePostal string, limit, offset int) (*models.CompanySearchResult, error) {
	conditions := []string{
		"e.etablissement_siege = 'true'",
		"e.code_postal_etablissement = $1",
	}
	args := []any{codePostal}
	cacheKey := fmt.Sprintf("sirene:v2:codepostal:%s", codePostal)
	criteria := models.CompanySearchCriteria{CodePostal: codePostal}

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
