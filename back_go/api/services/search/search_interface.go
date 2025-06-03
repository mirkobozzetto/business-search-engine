package search

import (
	"context"
	"csv-importer/api/models"
)

type SearchService interface {
	SearchInColumn(ctx context.Context, tableName, columnName, searchValue string, limit int) (*models.SearchResult, error)
	CountMatches(ctx context.Context, tableName, columnName, searchValue string) (*models.CountResult, error)
	SearchMultipleColumns(ctx context.Context, tableName string, columns []string, searchValue string, limit int) (*models.SearchResult, error)
	SearchNaceCode(ctx context.Context, searchValue string, limit int) (*models.NaceSearchResult, error)
}
