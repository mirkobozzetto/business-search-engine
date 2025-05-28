package utils

import (
	"csv-importer/api/models"
	"database/sql"
	"fmt"
)

func ScanRowsToMaps(rows *sql.Rows, columns []string) ([]map[string]any, error) {
	var data []map[string]any

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
	}

	return data, nil
}

func ScanRowsToMapsWithColumns(rows *sql.Rows) ([]map[string]any, []string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get columns: %v", err)
	}

	data, err := ScanRowsToMaps(rows, columns)
	if err != nil {
		return nil, nil, err
	}

	return data, columns, nil
}

func ScanColumnValues(rows *sql.Rows) ([]models.ColumnValue, error) {
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

	return values, nil
}

func BuildColumnStatsQuery(tableName, columnName string, limit int) (string, []any) {
	query := fmt.Sprintf(`
		SELECT %s, count(*) as count
		FROM %s
		WHERE %s IS NOT NULL AND %s != ''
		GROUP BY %s
		ORDER BY count DESC
		LIMIT $1
	`, columnName, tableName, columnName, columnName, columnName)

	return query, []any{limit}
}

func GetNullableValue(ns sql.NullString) any {
	if ns.Valid {
		return ns.String
	}
	return nil
}
