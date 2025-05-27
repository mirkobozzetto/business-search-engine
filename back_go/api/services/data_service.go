package services

import (
	"csv-importer/api/helpers"
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

	builder := &helpers.QueryBuilder{}
	builder.SetLimit(limit)
	rows, err := helpers.SafeQueryWithBuilder(s.db, tableName, nil, builder)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

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

	return &models.PreviewData{
		Table:   tableName,
		Columns: columns,
		Data:    data,
		Meta: models.Meta{
			Count: rowCount,
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

	query := fmt.Sprintf(`
		SELECT %s, count(*) as count
		FROM %s
		WHERE %s IS NOT NULL AND %s != ''
		GROUP BY %s
		ORDER BY count DESC
		LIMIT $1
	`, columnName, tableName, columnName, columnName, columnName)

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	var values []models.ColumnValue
	for rows.Next() {
		var value string
		var count int64
		if err := rows.Scan(&value, &count); err != nil {
			continue
		}
		values = append(values, models.ColumnValue{
			Value: value,
			Count: count,
		})
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
