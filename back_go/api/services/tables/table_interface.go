package tables

import (
	"context"
	"csv-importer/api/models"
)

type TableService interface {
	ListAllTables(ctx context.Context) ([]models.Table, error)
	GetTableInfo(ctx context.Context, tableName string) (*models.TableInfo, error)
	GetTableColumns(ctx context.Context, tableName string) ([]models.ColumnInfo, error)
	GetCompleteStructure(ctx context.Context) ([]models.TableStructure, error)
}
