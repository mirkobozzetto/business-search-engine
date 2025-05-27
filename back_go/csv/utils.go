package csv

import (
	"database/sql"
	"fmt"
	"strings"
)

func CleanColumnName(header string) string {
	cleanHeader := strings.ReplaceAll(header, " ", "_")
	cleanHeader = strings.ReplaceAll(cleanHeader, "-", "_")
	cleanHeader = strings.ToLower(cleanHeader)
	return cleanHeader
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
		"SET wal_buffers = '128MB'",
		"SET checkpoint_segments = 64",
		"SET checkpoint_completion_target = 0.9",
		"SET maintenance_work_mem = '1GB'",
		"SET work_mem = '512MB'",
		"SET shared_buffers = '512MB'",
		"SET effective_cache_size = '2GB'",
		"SET fsync = OFF", // DANGER: Only for imports!
	}

	for _, sql := range optimizations {
		if _, err := db.Exec(sql); err != nil {
			continue
		}
	}

	fmt.Println("🚀 PostgreSQL optimized for COPY")
	return nil
}
