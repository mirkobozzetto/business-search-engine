package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

func (s *companyService) SearchByDateCreation(ctx context.Context, fromDate, toDate string, limit int, offset int) (*models.CompanySearchResult, error) {
	conditions := []string{"e.etablissement_siege = 'true'"}
	var args []any
	argN := 1

	conditions = append(conditions, fmt.Sprintf("u.date_creation_unite_legale >= $%d", argN))
	args = append(args, fromDate)
	argN++

	if toDate != "" {
		conditions = append(conditions, fmt.Sprintf("u.date_creation_unite_legale <= $%d", argN))
		args = append(args, toDate)
		argN++
	}

	cacheKey := fmt.Sprintf("sirene:v2:datecreation:%s_%s", fromDate, toDate)
	criteria := models.CompanySearchCriteria{DateCreationFrom: fromDate, DateCreationTo: toDate}

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
