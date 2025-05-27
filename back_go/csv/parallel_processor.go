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

func ProcessCSVParallel(db *sql.DB, csvPath, tableName string) error {
	start := time.Now()
	numWorkers := 5
	chunkSize := 200000

	fmt.Printf("🚀 Starting PARALLEL import with %d workers\n", numWorkers)

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

	cleanHeaders, columns := PrepareHeaders(headers)

	OptimizeForBulkInsert(db)

	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	if _, err := db.Exec(dropSQL); err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))
	fmt.Printf("🏗️ Creating table: %s\n", tableName)

	if _, err := db.Exec(createSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	totalLines, err := processFileParallel(db, csvPath, tableName, cleanHeaders, numWorkers, chunkSize)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	linesPerSec := float64(totalLines) / elapsed.Seconds()

	fmt.Printf("🎉 PARALLEL IMPORT: %d lines in %.2f sec (%.0f lines/sec)\n",
		totalLines, elapsed.Seconds(), linesPerSec)

	return nil
}

func processFileParallel(db *sql.DB, csvPath, tableName string, headers []string, numWorkers, chunkSize int) (int, error) {
	chunks, err := CreateChunks(csvPath, chunkSize)
	if err != nil {
		return 0, err
	}

	pool := NewWorkerPool(db, tableName, headers, numWorkers)
	return pool.ProcessChunks(chunks)
}
