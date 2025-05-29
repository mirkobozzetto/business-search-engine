package data

import (
	"context"
	"csv-importer/api/models"
)

type DataService interface {
	PreviewTable(ctx context.Context, tableName string, limit int) (*models.PreviewData, error)
	GetColumnValues(ctx context.Context, tableName, columnName string, limit int) (*models.ColumnValues, error)
}
