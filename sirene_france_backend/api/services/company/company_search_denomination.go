package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
	"strings"
)

func (s *companyService) SearchByDenomination(ctx context.Context, query string, limit, offset int) (*models.CompanySearchResult, error) {
	words := strings.Fields(strings.TrimSpace(query))
	if len(words) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{Denomination: query},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	conditions := []string{"e.etablissement_siege = 'true'"}
	var args []any
	argN := 1

	for _, word := range words {
		conditions = append(conditions, fmt.Sprintf("unaccent(u.denomination_unite_legale) ILIKE unaccent($%d)", argN))
		args = append(args, "%"+word+"%")
		argN++
	}

	cacheKey := fmt.Sprintf("sirene:v2:denomination:%s", query)
	criteria := models.CompanySearchCriteria{Denomination: query}

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
