package utils

import (
	"fmt"
	"strings"
)

type SearchBuilder struct {
	columns []string
	args    []any
	conditions []string
}

func NewSearchBuilder(columns []string) *SearchBuilder {
	return &SearchBuilder{
		columns: columns,
		args:    make([]any, 0),
		conditions: make([]string, 0),
	}
}

func (sb *SearchBuilder) BuildMultiWordSearch(searchValue string) (string, []any) {
	if searchValue == "" {
		return "", sb.args
	}

	words := strings.Fields(strings.TrimSpace(searchValue))
	if len(words) == 0 {
		return "", sb.args
	}

	var columnConditions []string

	for _, column := range sb.columns {
		var wordConditions []string

		for _, word := range words {
			sb.args = append(sb.args, "%"+word+"%")
			wordConditions = append(wordConditions, fmt.Sprintf("%s ILIKE $%d", column, len(sb.args)))
		}

		if len(wordConditions) > 1 {
			columnConditions = append(columnConditions, "("+strings.Join(wordConditions, " AND ")+")")
		} else {
			columnConditions = append(columnConditions, wordConditions[0])
		}
	}

	whereClause := strings.Join(columnConditions, " OR ")

	if len(columnConditions) > 1 {
		whereClause = "(" + whereClause + ")"
	}

	return whereClause, sb.args
}

func (sb *SearchBuilder) BuildSimpleSearch(searchValue string) (string, []any) {
	if searchValue == "" {
		return "", sb.args
	}

	var columnConditions []string

	for _, column := range sb.columns {
		sb.args = append(sb.args, "%"+searchValue+"%")
		columnConditions = append(columnConditions, fmt.Sprintf("%s ILIKE $%d", column, len(sb.args)))
	}

	whereClause := strings.Join(columnConditions, " OR ")

	if len(columnConditions) > 1 {
		whereClause = "(" + whereClause + ")"
	}

	return whereClause, sb.args
}

func ParseOptionalLimit(limitStr string, defaultLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}

	limit := 0
	if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil || limit <= 0 {
		return defaultLimit
	}

	return limit
}

func BuildNaceCodeQuery(searchValue string, limit int) (string, []any) {
	columns := []string{"activités", "libellé_fr", "omschrijving_nl"}
	builder := NewSearchBuilder(columns)

	selectClause := "SELECT nacecode, activités, libellé_fr, omschrijving_nl FROM nacecode"

	if searchValue == "" {
		if limit > 0 {
			return fmt.Sprintf("%s ORDER BY nacecode LIMIT $1", selectClause), []any{limit}
		}
		return fmt.Sprintf("%s ORDER BY nacecode", selectClause), []any{}
	}

	whereClause, args := builder.BuildMultiWordSearch(searchValue)

	if limit > 0 {
		args = append(args, limit)
		return fmt.Sprintf("%s WHERE %s ORDER BY nacecode LIMIT $%d",
			selectClause, whereClause, len(args)), args
	}

	return fmt.Sprintf("%s WHERE %s ORDER BY nacecode", selectClause, whereClause), args
}
