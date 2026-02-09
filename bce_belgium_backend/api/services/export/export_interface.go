package export

import (
	"context"
	"csv-importer/api/models"
)

type ExportOptions struct {
	TableName   string
	ColumnName  string
	SearchValue string
	Limit       int
	Format      string
}

type ExportResult struct {
	Data     []map[string]any `json:"data,omitempty"`
	Columns  []string         `json:"columns"`
	RowCount int              `json:"row_count"`
	Meta     models.Meta      `json:"meta"`
}

type ExportService interface {
	PrepareExportData(ctx context.Context, opts ExportOptions) (*ExportResult, error)
	GenerateFilename(opts ExportOptions) string
	ValidateExportOptions(opts ExportOptions) error
}
