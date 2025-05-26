package query

import (
	"database/sql"
	"fmt"
	"strings"
)

func ListTables(db *sql.DB) error {
	query := `
		SELECT
			table_name,
			pg_size_pretty(pg_total_relation_size(table_name::regclass)) as size,
			(SELECT count(*) FROM information_schema.columns WHERE table_name = t.table_name) as columns
		FROM information_schema.tables t
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		ORDER BY pg_total_relation_size(table_name::regclass) DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying tables: %v", err)
	}
	defer rows.Close()

	fmt.Println("\n📊 TABLES CREATED:")
	fmt.Printf("%-20s %-10s %-8s\n", "TABLE", "SIZE", "COLUMNS")
	fmt.Println(strings.Repeat("-", 40))

	for rows.Next() {
		var tableName, size string
		var columns int

		if err := rows.Scan(&tableName, &size, &columns); err != nil {
			continue
		}

		fmt.Printf("%-20s %-10s %-8d\n", tableName, size, columns)
	}

	return nil
}

func CountRows(db *sql.DB, tableName string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT count(*) FROM %s", tableName)

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting rows in %s: %v", tableName, err)
	}

	return count, nil
}

func ShowTableInfo(db *sql.DB, tableName string) error {
	// Get column info
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return fmt.Errorf("error getting table info: %v", err)
	}
	defer rows.Close()

	fmt.Printf("\n🔍 TABLE: %s\n", strings.ToUpper(tableName))
	fmt.Printf("%-25s %-15s %-10s\n", "COLUMN", "TYPE", "NULLABLE")
	fmt.Println(strings.Repeat("-", 55))

	for rows.Next() {
		var columnName, dataType, nullable string

		if err := rows.Scan(&columnName, &dataType, &nullable); err != nil {
			continue
		}

		fmt.Printf("%-25s %-15s %-10s\n", columnName, dataType, nullable)
	}

	// Get row count
	count, err := CountRows(db, tableName)
	if err == nil {
		fmt.Printf("\n📊 Total rows: %d\n", count)
	}

	return nil
}

func PreviewTable(db *sql.DB, tableName string, limit int) error {
	if limit <= 0 {
		limit = 5
	}

	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, limit)

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error previewing table: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error getting columns: %v", err)
	}

	fmt.Printf("\n👀 PREVIEW: %s (first %d rows)\n", strings.ToUpper(tableName), limit)

	// Print headers
	for i, col := range columns {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Printf("%-15s", col)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", len(columns)*18))

	// Print rows
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		for i, val := range values {
			if i > 0 {
				fmt.Print(" | ")
			}

			str := ""
			if val.Valid {
				str = val.String
				if len(str) > 15 {
					str = str[:12] + "..."
				}
			}
			fmt.Printf("%-15s", str)
		}
		fmt.Println()
	}

	return nil
}

func ShowStats(db *sql.DB) error {
	fmt.Println("\n🚀 DATABASE STATISTICS:")

	// Get total rows across all tables
	query := `
		SELECT
			schemaname,
			tablename,
			n_tup_ins as inserts,
			n_tup_upd as updates,
			n_tup_del as deletes,
			n_live_tup as live_rows
		FROM pg_stat_user_tables
		ORDER BY n_live_tup DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error getting stats: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-20s %-12s %-12s %-12s %-12s\n", "TABLE", "LIVE_ROWS", "INSERTS", "UPDATES", "DELETES")
	fmt.Println(strings.Repeat("-", 75))

	var totalRows int64
	for rows.Next() {
		var schema, table string
		var inserts, updates, deletes, liveRows int64

		if err := rows.Scan(&schema, &table, &inserts, &updates, &deletes, &liveRows); err != nil {
			continue
		}

		fmt.Printf("%-20s %-12d %-12d %-12d %-12d\n", table, liveRows, inserts, updates, deletes)
		totalRows += liveRows
	}

	fmt.Println(strings.Repeat("-", 75))
	fmt.Printf("🎯 TOTAL ROWS: %d (%.1fM)\n", totalRows, float64(totalRows)/1000000)

	return nil
}
