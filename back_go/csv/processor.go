package csv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ProcessCSVOptimized(db *sql.DB, csvPath, tableName string) error {
	start := time.Now()

	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("impossible to open %s: %v", csvPath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading header: %v", err)
	}

	fmt.Printf("📄 CSV: %s\n", filepath.Base(csvPath))
	fmt.Printf("📊 Columns: %v\n", headers)

	// Optimize PostgreSQL settings for bulk insert
	optimizeForBulkInsert(db)

	// Drop table
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	if _, err := db.Exec(dropSQL); err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	// Create table with optimizations
	var columns []string
	cleanHeaders := make([]string, len(headers))
	for i, header := range headers {
		cleanHeader := strings.ReplaceAll(header, " ", "_")
		cleanHeader = strings.ReplaceAll(cleanHeader, "-", "_")
		cleanHeader = strings.ToLower(cleanHeader)
		cleanHeaders[i] = cleanHeader
		columns = append(columns, cleanHeader+" TEXT") // TEXT instead of VARCHAR(255)
	}

	createSQL := fmt.Sprintf("CREATE UNLOGGED TABLE %s (%s)", tableName, strings.Join(columns, ", "))
	fmt.Printf("🏗️ Creating UNLOGGED table: %s\n", tableName)

	if _, err := db.Exec(createSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Use optimized batch insert (easier than COPY for now)
	lineCount, err := optimizedBatchInsert(db, reader, tableName, cleanHeaders)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	linesPerSec := float64(lineCount) / elapsed.Seconds()

	fmt.Printf("✅ Total: %d lines in %.2f sec (%.0f lines/sec)\n",
		lineCount, elapsed.Seconds(), linesPerSec)
	return nil
}

func optimizedBatchInsert(db *sql.DB, reader *csv.Reader, tableName string, headers []string) (int, error) {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	lineCount := 0
	// Calculate optimal batch size based on PostgreSQL parameter limit
	maxParams := 65000 // Slightly under the 65535 limit
	optimalBatchSize := maxParams / len(headers)
	if optimalBatchSize > 15000 {
		optimalBatchSize = 15000 // Cap at reasonable size
	}

	batchSize := optimalBatchSize
	fmt.Printf("🎯 Optimal batch size: %d lines (%d params per batch)\n",
		batchSize, batchSize*len(headers))
	batch := make([][]any, 0, batchSize)
	startTime := time.Now()

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return 0, fmt.Errorf("error reading line %d: %v", lineCount+1, err)
		}

		// Add to batch
		values := make([]any, len(record))
		for i, v := range record {
			values[i] = v
		}
		batch = append(batch, values)
		lineCount++

		// Insert batch when full
		if len(batch) >= batchSize {
			if err := insertBigBatch(tx, tableName, len(headers), batch); err != nil {
				return 0, fmt.Errorf("error inserting batch at line %d: %v", lineCount, err)
			}

			elapsed := time.Since(startTime)
			linesPerSec := float64(lineCount) / elapsed.Seconds()
			fmt.Printf("📈 Processed: %d lines (%.0f lines/sec)\n", lineCount, linesPerSec)

			batch = batch[:0] // Reset batch
		}
	}

	// Insert remaining batch
	if len(batch) > 0 {
		if err := insertBigBatch(tx, tableName, len(headers), batch); err != nil {
			return 0, fmt.Errorf("error inserting final batch: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return lineCount, nil
}

func insertBigBatch(tx *sql.Tx, tableName string, numCols int, batch [][]any) error {
	if len(batch) == 0 {
		return nil
	}

	// Build INSERT with multiple VALUES - bigger batch
	var valuePlaceholders []string
	var allValues []any

	for i, row := range batch {
		var placeholders []string
		for j := range numCols {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i*numCols+j+1))
			allValues = append(allValues, row[j])
		}
		valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ", ")+")")
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s VALUES %s",
		tableName, strings.Join(valuePlaceholders, ", "))

	_, err := tx.Exec(insertSQL, allValues...)
	return err
}

func optimizeForBulkInsert(db *sql.DB) error {
	optimizations := []string{
		"SET synchronous_commit = OFF",
		"SET wal_buffers = '64MB'",
		"SET checkpoint_segments = 32",
		"SET checkpoint_completion_target = 0.9",
		"SET maintenance_work_mem = '512MB'",
		"SET work_mem = '256MB'",
	}

	for _, sql := range optimizations {
		if _, err := db.Exec(sql); err != nil {
			// Ignore errors for settings that might not exist
			continue
		}
	}

	fmt.Println("🚀 PostgreSQL optimized for bulk insert")
	return nil
}

// Keep the old function for compatibility
func ProcessCSV(db *sql.DB, csvPath, tableName string) error {
	return ProcessCSVOptimized(db, csvPath, tableName)
}
