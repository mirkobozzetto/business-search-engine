package utils

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

var ValidOperators = []string{
	"=", "!=", "<>", "<", ">", "<=", ">=",
	"LIKE", "ILIKE", "IN", "NOT IN", "IS NULL", "IS NOT NULL",
}

type QueryBuilder struct {
	conditions []string
	args       []any
	limit      *int
}

func (qb *QueryBuilder) AddCondition(column string, operator string, value any) error {
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

func JoinColumns(columns []string) string {
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

func BuildSafeQuery(tableName string, columns []string, builder *QueryBuilder) (string, []any) {
	query := fmt.Sprintf("SELECT %s FROM %s", JoinColumns(columns), tableName)

	whereClause, args := builder.BuildWhere()
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	if builder.limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *builder.limit)
	}

	return query, args
}

func ExecuteSafeQuery(db *sql.DB, query string, args []any) (*sql.Rows, error) {
	return db.Query(query, args...)
}
