package data

import (
	"context"
	"csv-importer/api/helpers"
	"csv-importer/api/helpers/utils"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

type dataService struct {
	db *sql.DB
}

// NewDataService creates a new DataService implementation
func NewDataService(db *sql.DB) DataService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	return &dataService{
		db: db,
	}
}

func (s *dataService) PreviewTable(ctx context.Context, tableName string, limit int) (*models.PreviewData, error) {
	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	if limit <= 0 || limit > 100 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 100")
	}

	builder := &utils.QueryBuilder{}
	builder.SetLimit(limit)

	rows, err := helpers.SafeQueryWithBuilder(s.db, tableName, nil, builder)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	data, columns, err := utils.ScanRowsToMapsWithColumns(rows)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return &models.PreviewData{
		Table:   tableName,
		Columns: columns,
		Data:    data,
		Meta: models.Meta{
			Count: len(data),
			Limit: limit,
		},
	}, nil
}

func (s *dataService) GetColumnValues(ctx context.Context, tableName, columnName string, limit int) (*models.ColumnValues, error) {
	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	if err := helpers.ValidateColumnExists(s.db, tableName, columnName); err != nil {
		return nil, fmt.Errorf("invalid column: %w", err)
	}

	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	query, args := utils.BuildColumnStatsQuery(tableName, columnName, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	values, err := utils.ScanColumnValues(rows)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return &models.ColumnValues{
		Table:  tableName,
		Column: columnName,
		Values: values,
		Meta: models.Meta{
			Count: len(values),
			Limit: limit,
		},
	}, nil
}
