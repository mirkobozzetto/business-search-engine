package company

import (
	"context"
	"csv-importer/api/models"
)

type CompanyService interface {
	SearchByNaceCode(ctx context.Context, naceCode string, limit int) (*models.CompanySearchResult, error)
}
