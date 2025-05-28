package services

import (
	"csv-importer/api/helpers"
	"csv-importer/api/helpers/utils"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
)

type DataService struct {
	db *sql.DB
}

func NewDataService(db *sql.DB) *DataService {
	return &DataService{db: db}
}

func (s *DataService) PreviewTable(tableName string, limit int) (*models.PreviewData, error) {
	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 100 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 100")
	}

	builder := &utils.QueryBuilder{}
	builder.SetLimit(limit)
	rows, err := helpers.SafeQueryWithBuilder(s.db, tableName, nil, builder)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	data, columns, err := utils.ScanRowsToMapsWithColumns(rows)
	if err != nil {
		return nil, err
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

func (s *DataService) GetColumnValues(tableName, columnName string, limit int) (*models.ColumnValues, error) {
	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, err
	}

	if err := helpers.ValidateColumnExists(s.db, tableName, columnName); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	query, args := utils.BuildColumnStatsQuery(tableName, columnName, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	values, err := utils.ScanColumnValues(rows)
	if err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
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
