package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

func (s *companyService) SearchByEtatAdministratif(ctx context.Context, etat string, limit int, offset int) (*models.CompanySearchResult, error) {
	conditions := []string{
		"e.etablissement_siege = 'true'",
		"u.etat_administratif_unite_legale = $1",
	}
	args := []any{etat}
	cacheKey := fmt.Sprintf("sirene:v2:etatadmin:%s", etat)
	criteria := models.CompanySearchCriteria{EtatAdministratif: etat}

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
