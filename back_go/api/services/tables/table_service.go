package tables

import (
	"context"
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

type tableService struct {
	db *sql.DB
}

func NewTableService(db *sql.DB) TableService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	return &tableService{
		db: db,
	}
}

func (s *tableService) ListAllTables(ctx context.Context) ([]models.Table, error) {
	query := `
		SELECT table_name,
		       (SELECT count(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
		FROM information_schema.tables t
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var tableName string
		var columnCount int
		if err := rows.Scan(&tableName, &columnCount); err != nil {
			continue
		}

		rowCount, err := s.getTableRowCount(tableName)
		if err != nil {
			continue
		}

		tables = append(tables, models.Table{
			Name:    tableName,
			Rows:    rowCount,
			Columns: columnCount,
		})
	}

	return tables, nil
}

func (s *tableService) GetTableInfo(ctx context.Context, tableName string) (*models.TableInfo, error) {
	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	colCount, err := s.getTableColumnCount(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get column count: %w", err)
	}

	rowCount, err := s.getTableRowCount(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get row count: %w", err)
	}

	fields, err := s.getTableFields(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table fields: %w", err)
	}

	return &models.TableInfo{
		Table:   tableName,
		Rows:    rowCount,
		Columns: colCount,
		Fields:  fields,
	}, nil
}

func (s *tableService) GetTableColumns(ctx context.Context, tableName string) ([]models.ColumnInfo, error) {
	if err := helpers.ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := s.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	defer rows.Close()

	var columns []models.ColumnInfo
	for rows.Next() {
		var name, dataType, nullable string
		if err := rows.Scan(&name, &dataType, &nullable); err != nil {
			continue
		}

		columns = append(columns, models.ColumnInfo{
			Name:     name,
			Type:     dataType,
			Nullable: nullable == "YES",
		})
	}

	return columns, nil
}

func (s *tableService) GetCompleteStructure(ctx context.Context) ([]models.TableStructure, error) {
	tablesQuery := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := s.db.Query(tablesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var structures []models.TableStructure
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}

		rowCount, err := s.getTableRowCount(tableName)
		if err != nil {
			continue
		}

		columns, err := s.GetTableColumns(context.Background(), tableName)
		if err != nil {
			continue
		}

		structures = append(structures, models.TableStructure{
			Name:    tableName,
			Rows:    rowCount,
			Columns: columns,
		})
	}

	return structures, nil
}

func (s *tableService) getTableRowCount(tableName string) (int64, error) {
	var rowCount int64
	query := fmt.Sprintf("SELECT count(*) FROM %s", tableName)
	err := s.db.QueryRow(query).Scan(&rowCount)
	return rowCount, err
}

func (s *tableService) getTableColumnCount(tableName string) (int, error) {
	var colCount int
	query := `SELECT count(*) FROM information_schema.columns WHERE table_name = $1`
	err := s.db.QueryRow(query, tableName).Scan(&colCount)
	return colCount, err
}

func (s *tableService) getTableFields(tableName string) ([]string, error) {
	query := `SELECT column_name FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position`
	rows, err := s.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}
	defer rows.Close()

	var fields []string
	for rows.Next() {
		var field string
		if err := rows.Scan(&field); err != nil {
			continue
		}
		fields = append(fields, field)
	}

	return fields, nil
}
