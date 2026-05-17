package csv

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ProcessCSVBlazingFast(db *sql.DB, csvPath, tableName string) error {
	start := time.Now()

	fmt.Printf("🔥 ULTRA FAST streaming import\n")

	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("impossible to open %s: %v", csvPath, err)
	}
	defer func() { _ = file.Close() }()

	bufferedReader := bufio.NewReaderSize(file, 2*1024*1024)
	reader := csv.NewReader(bufferedReader)
	reader.Comma = ','

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading header: %v", err)
	}

	fmt.Printf("📄 CSV: %s\n", filepath.Base(csvPath))
	fmt.Printf("📊 Columns: %v\n", headers)

	cleanHeaders, columns := PrepareHeaders(headers)

	if err := setupTable(db, tableName, columns); err != nil {
		return err
	}

	totalLines, err := ProcessPipelineParallel(db, csvPath, tableName, cleanHeaders)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	linesPerSec := float64(totalLines) / elapsed.Seconds()

	fmt.Printf("🔥 ULTRA FAST: %d lines in %.2f sec (%.0f lines/sec)\n",
		totalLines, elapsed.Seconds(), linesPerSec)

	return nil
}

func setupTable(db *sql.DB, tableName string, columns []string) error {
	_ = OptimizeForBulkInsert(db)

	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	if _, err := db.Exec(dropSQL); err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))
	fmt.Printf("🏗️ Creating table: %s\n", tableName)

	if _, err := db.Exec(createSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}
