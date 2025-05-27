package helpers

import (
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var ValidTables = []string{
	"activity", "address", "branch", "code", "contact",
	"denomination", "enterprise", "establishment", "meta", "nacecode",
}

var ValidOperators = []string{
	"=", "!=", "<>", "<", ">", "<=", ">=",
	"LIKE", "ILIKE", "IN", "NOT IN", "IS NULL", "IS NOT NULL",
}

func ValidateTableName(tableName string) error {
	if slices.Contains(ValidTables, tableName) {
		return nil
	}
	return fmt.Errorf("invalid table name: %s", tableName)
}

func ValidateColumnExists(db *sql.DB, tableName, columnName string) error {
	var count int
	err := db.QueryRow(`
		SELECT count(*)
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2
	`, tableName, columnName).Scan(&count)

	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("column %s not found in table %s", columnName, tableName)
	}
	return nil
}

func SafeQuery(db *sql.DB, tableName string, columns []string) (*sql.Rows, error) {
	if err := ValidateTableName(tableName); err != nil {
		return nil, err
	}

	for _, col := range columns {
		if err := ValidateIdentifier(col); err != nil {
			return nil, err
		}
		if err := ValidateColumnExists(db, tableName, col); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", joinColumns(columns), tableName)
	return db.Query(query)
}

func joinColumns(columns []string) string {
	if len(columns) == 0 {
		return "*"
	}

	result := ""
	for i, col := range columns {
		if i > 0 {
			result += ", "
		}
		result += col
	}
	return result
}

type QueryBuilder struct {
	conditions []string
	args       []any
	limit      *int
}

func (qb *QueryBuilder) AddCondition(column string, operator string, value any) error {
	if err := ValidateIdentifier(column); err != nil {
		return err
	}

	if !slices.Contains(ValidOperators, strings.ToUpper(operator)) {
		return fmt.Errorf("invalid operator: %s", operator)
	}

	qb.conditions = append(qb.conditions, fmt.Sprintf("%s %s $%d", column, operator, len(qb.args)+1))
	qb.args = append(qb.args, value)
	return nil
}

func (qb *QueryBuilder) SetLimit(limit int) {
	qb.limit = &limit
}

func (qb *QueryBuilder) BuildWhere() (string, []any) {
	if len(qb.conditions) == 0 {
		return "", qb.args
	}
	return strings.Join(qb.conditions, " AND "), qb.args
}

func SafeQueryWithBuilder(db *sql.DB, tableName string, columns []string, builder *QueryBuilder) (*sql.Rows, error) {
	if err := ValidateTableName(tableName); err != nil {
		return nil, err
	}

	for _, col := range columns {
		if err := ValidateIdentifier(col); err != nil {
			return nil, err
		}
		if err := ValidateColumnExists(db, tableName, col); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", joinColumns(columns), tableName)

	whereClause, args := builder.BuildWhere()
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	if builder.limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *builder.limit)
	}

	return db.Query(query, args...)
}

func ValidateIdentifier(name string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	if !matched {
		return fmt.Errorf("invalid identifier: %s", name)
	}
	return nil
}
