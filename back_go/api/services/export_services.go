package services

import (
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"strings"
)

type ExportService struct {
	db *sql.DB
}

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

func NewExportService(db *sql.DB) *ExportService {
	return &ExportService{db: db}
}

func (s *ExportService) ValidateExportOptions(opts ExportOptions) error {
	if err := helpers.ValidateTableName(opts.TableName); err != nil {
		return err
	}

	if opts.Limit <= 0 || opts.Limit > 100000 {
		return fmt.Errorf("invalid limit: must be between 1 and 100000")
	}

	if opts.ColumnName != "" {
		if err := helpers.ValidateColumnExists(s.db, opts.TableName, opts.ColumnName); err != nil {
			return err
		}
	}

	if opts.Format != "" && opts.Format != "csv" && opts.Format != "json" {
		return fmt.Errorf("invalid format: must be 'csv' or 'json'")
	}

	return nil
}

func (s *ExportService) GetTableColumns(tableName string) ([]string, error) {
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`
	rows, err := s.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %v", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			continue
		}
		columns = append(columns, column)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("table not found or has no columns")
	}

	return columns, nil
}

func (s *ExportService) BuildExportQuery(opts ExportOptions, columns []string) (string, []any) {
	var query string
	var args []any

	if opts.ColumnName != "" && opts.SearchValue != "" {
		query = fmt.Sprintf(`
			SELECT %s
			FROM %s
			WHERE %s ILIKE $1
			LIMIT $2
		`, strings.Join(columns, ","), opts.TableName, opts.ColumnName)
		args = []any{"%" + opts.SearchValue + "%", opts.Limit}
	} else {
		query = fmt.Sprintf(`
			SELECT %s
			FROM %s
			LIMIT $1
		`, strings.Join(columns, ","), opts.TableName)
		args = []any{opts.Limit}
	}

	return query, args
}

func (s *ExportService) ExecuteExportQuery(query string, args []any, columns []string) (*ExportResult, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	var data []map[string]any
	rowCount := 0

	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		row := make(map[string]any)
		for i, col := range columns {
			if values[i].Valid {
				row[col] = values[i].String
			} else {
				row[col] = nil
			}
		}
		data = append(data, row)
		rowCount++
	}

	return &ExportResult{
		Data:     data,
		Columns:  columns,
		RowCount: rowCount,
		Meta: models.Meta{
			Count: rowCount,
			Limit: len(args) - 1, // Last arg is usually limit
		},
	}, nil
}

func (s *ExportService) PrepareExportData(opts ExportOptions) (*ExportResult, error) {
	if err := s.ValidateExportOptions(opts); err != nil {
		return nil, err
	}

	columns, err := s.GetTableColumns(opts.TableName)
	if err != nil {
		return nil, err
	}

	query, args := s.BuildExportQuery(opts, columns)

	return s.ExecuteExportQuery(query, args, columns)
}

func (s *ExportService) GenerateFilename(opts ExportOptions) string {
	if opts.SearchValue != "" {
		return fmt.Sprintf("%s_%s_export", opts.TableName, opts.SearchValue)
	}
	return fmt.Sprintf("%s_export", opts.TableName)
}
