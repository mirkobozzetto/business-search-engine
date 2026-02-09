package export

import (
	"context"
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type exportService struct {
	db *sql.DB
}

func NewExportService(db *sql.DB) ExportService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	return &exportService{
		db: db,
	}
}

func (s *exportService) ValidateExportOptions(opts ExportOptions) error {
	if err := helpers.ValidateTableName(opts.TableName); err != nil {
		return fmt.Errorf("invalid table name: %w", err)
	}

	if opts.Limit <= 0 || opts.Limit > 100000 {
		return fmt.Errorf("invalid limit: must be between 1 and 100000")
	}

	if opts.ColumnName != "" {
		if err := helpers.ValidateColumnExists(s.db, opts.TableName, opts.ColumnName); err != nil {
			return fmt.Errorf("invalid column: %w", err)
		}
	}

	if opts.Format != "" && opts.Format != "csv" && opts.Format != "json" {
		return fmt.Errorf("invalid format: must be 'csv' or 'json'")
	}

	return nil
}

func (s *exportService) getTableColumns(tableName string) ([]string, error) {
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`
	rows, err := s.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
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

func (s *exportService) buildExportQuery(opts ExportOptions, columns []string) (string, []any) {
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

func (s *exportService) executeExportQuery(query string, args []any, columns []string, limit int) (*ExportResult, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
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
			Limit: limit,
		},
	}, nil
}

func (s *exportService) PrepareExportData(ctx context.Context, opts ExportOptions) (*ExportResult, error) {
	if err := s.ValidateExportOptions(opts); err != nil {
		return nil, err
	}

	columns, err := s.getTableColumns(opts.TableName)
	if err != nil {
		return nil, err
	}

	query, args := s.buildExportQuery(opts, columns)

	return s.executeExportQuery(query, args, columns, opts.Limit)
}

func (s *exportService) GenerateFilename(opts ExportOptions) string {
	if opts.SearchValue != "" {
		return fmt.Sprintf("%s_%s_export", opts.TableName, opts.SearchValue)
	}
	return fmt.Sprintf("%s_export", opts.TableName)
}
