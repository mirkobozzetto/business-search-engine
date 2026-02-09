package search

import (
	"context"
	"csv-importer/api/helpers"
	helperutils "csv-importer/api/helpers/utils"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type searchService struct {
	db *sql.DB
}

func NewSearchService(db *sql.DB) SearchService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	return &searchService{
		db: db,
	}
}

func (s *searchService) SearchInColumn(ctx context.Context, tableName, columnName, searchValue string, limit int) (*models.SearchResult, error) {
	if searchValue == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	if err := helpers.ValidateColumnExists(s.db, tableName, columnName); err != nil {
		return nil, fmt.Errorf("invalid column: %w", err)
	}

	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM %s
		WHERE %s ILIKE $1
		ORDER BY %s
		LIMIT $2
	`, columnName, tableName, columnName, columnName)

	rows, err := s.db.Query(query, "%"+searchValue+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var value sql.NullString
		if err := rows.Scan(&value); err != nil {
			continue
		}
		if value.Valid {
			results = append(results, value.String)
		}
	}

	return &models.SearchResult{
		Table:   tableName,
		Column:  columnName,
		Query:   searchValue,
		Results: results,
		Meta: models.Meta{
			Count: len(results),
			Limit: limit,
		},
	}, nil
}

func (s *searchService) CountMatches(ctx context.Context, tableName, columnName, searchValue string) (*models.CountResult, error) {
	if searchValue == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	if err := helpers.ValidateColumnExists(s.db, tableName, columnName); err != nil {
		return nil, fmt.Errorf("invalid column: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE %s ILIKE $1
	`, tableName, columnName)

	var count int64
	err := s.db.QueryRow(query, "%"+searchValue+"%").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &models.CountResult{
		Table:  tableName,
		Column: columnName,
		Query:  searchValue,
		Count:  count,
	}, nil
}

func (s *searchService) SearchMultipleColumns(ctx context.Context, tableName string, columns []string, searchValue string, limit int) (*models.SearchResult, error) {
	if searchValue == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	for _, col := range columns {
		if err := helpers.ValidateColumnExists(s.db, tableName, col); err != nil {
			return nil, fmt.Errorf("invalid column %s: %w", col, err)
		}
	}

	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	// Build WHERE clause for multiple columns
	var whereConditions []string
	for _, col := range columns {
		whereConditions = append(whereConditions, fmt.Sprintf("%s ILIKE $1", col))
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM %s
		WHERE %s
		LIMIT $2
	`, strings.Join(columns, ", "), tableName, strings.Join(whereConditions, " OR "))

	rows, err := s.db.Query(query, "%"+searchValue+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		var rowValues []string
		for _, val := range values {
			if val.Valid {
				rowValues = append(rowValues, val.String)
			}
		}
		if len(rowValues) > 0 {
			results = append(results, strings.Join(rowValues, " | "))
		}
	}

	return &models.SearchResult{
		Table:   tableName,
		Column:  strings.Join(columns, ","),
		Query:   searchValue,
		Results: results,
		Meta: models.Meta{
			Count: len(results),
			Limit: limit,
		},
	}, nil
}

func (s *searchService) SearchNaceCode(ctx context.Context, searchValue string, limit int) (*models.NaceSearchResult, error) {
	query, args := buildNaceCodeQuery(searchValue, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	columns := []string{"nacecode", "activités", "libellé_fr", "omschrijving_nl"}
	data, err := helperutils.ScanRowsToMaps(rows, columns)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return &models.NaceSearchResult{
		Query:   searchValue,
		Results: data,
		Meta: models.Meta{
			Count: len(data),
			Limit: limit,
			Total: len(data),
		},
	}, nil
}
