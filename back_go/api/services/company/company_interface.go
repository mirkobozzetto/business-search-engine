package company

import (
	"context"
	"csv-importer/api/models"
)

type CompanyService interface {
	SearchByNaceCode(ctx context.Context, naceCode string, limit int) (*models.CompanySearchResult, error)
	SearchByDenomination(ctx context.Context, query string, limit int) (*models.CompanySearchResult, error)
	SearchByZipcode(ctx context.Context, zipcode string, limit int) (*models.CompanySearchResult, error)
	SearchMultiCriteria(ctx context.Context, criteria models.CompanySearchCriteria, limit int) (*models.CompanySearchResult, error)
}
