package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

func (s *companyService) SearchByCommune(ctx context.Context, commune string, limit, offset int) (*models.CompanySearchResult, error) {
	conditions := []string{
		"e.etablissement_siege = 'true'",
		"unaccent(e.libelle_commune_etablissement) ILIKE unaccent($1)",
	}
	args := []any{"%" + commune + "%"}
	cacheKey := fmt.Sprintf("sirene:v2:commune:%s", commune)
	criteria := models.CompanySearchCriteria{Commune: commune}

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
