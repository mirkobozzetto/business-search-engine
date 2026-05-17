package csv

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

var camelCaseRe = regexp.MustCompile("([a-z0-9])([A-Z])")

func CleanColumnName(header string) string {
	snake := camelCaseRe.ReplaceAllString(header, "${1}_${2}")
	snake = strings.ReplaceAll(snake, " ", "_")
	snake = strings.ReplaceAll(snake, "-", "_")
	return strings.ToLower(snake)
}

func PrepareHeaders(headers []string) ([]string, []string) {
	cleanHeaders := make([]string, len(headers))
	var columns []string

	for i, header := range headers {
		cleanHeader := CleanColumnName(header)
		cleanHeaders[i] = cleanHeader
		columns = append(columns, cleanHeader+" TEXT")
	}

	return cleanHeaders, columns
}

func OptimizeForBulkInsert(db *sql.DB) error {
	optimizations := []string{
		"SET synchronous_commit = OFF",
		"SET maintenance_work_mem = '1GB'",
		"SET work_mem = '512MB'",
	}

	for _, sql := range optimizations {
		if _, err := db.Exec(sql); err != nil {
			continue
		}
	}

	fmt.Println("PostgreSQL optimized for COPY")
	return nil
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
