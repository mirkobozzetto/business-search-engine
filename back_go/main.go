package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Pas de fichier .env trouvé, utilisation variables système")
	}
	dbHost := getEnv("DB_HOST", "")
	dbPort := getEnv("DB_PORT", "")
	dbUser := getEnv("POSTGRES_USER", "")
	dbPassword := getEnv("POSTGRES_PASSWORD", "")
	dbName := getEnv("POSTGRES_DB", "")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("DB not connected: %v", err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("DB not connected: %v", err))
	}

	fmt.Println("✅ DB connected")

	csvPath := "../bce_mai_2025/activity.csv"
	tableName := "activity"

	if err := processCSV(db, csvPath, tableName); err != nil {
		panic(fmt.Sprintf("Error processing CSV: %v", err))
	}

	fmt.Println("🎉 Import terminé!")
}

func processCSV(db *sql.DB, csvPath, tableName string) error {
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

	// Drop table
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	if _, err := db.Exec(dropSQL); err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	// Create table
	var columns []string
	for _, header := range headers {
		cleanHeader := strings.ReplaceAll(header, " ", "_")
		cleanHeader = strings.ReplaceAll(cleanHeader, "-", "_")
		cleanHeader = strings.ToLower(cleanHeader)
		columns = append(columns, cleanHeader+" VARCHAR(255)")
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))
	fmt.Printf("🏗️ Creating table: %s\n", tableName)

	if _, err := db.Exec(createSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Batch insert setup
	lineCount := 0
	batchSize := 5000 // Plus gros batch
	batch := make([][]any, 0, batchSize)

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading line %d: %v", lineCount+1, err)
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
			if err := insertBatch(tx, tableName, len(headers), batch); err != nil {
				return fmt.Errorf("error inserting batch at line %d: %v", lineCount, err)
			}
			fmt.Printf("📈 Processed: %d lines\n", lineCount)
			batch = batch[:0] // Reset batch
		}
	}

	// Insert remaining batch
	if len(batch) > 0 {
		if err := insertBatch(tx, tableName, len(headers), batch); err != nil {
			return fmt.Errorf("error inserting final batch: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	fmt.Printf("✅ Total inserted: %d lines\n", lineCount)
	return nil
}

func insertBatch(tx *sql.Tx, tableName string, numCols int, batch [][]any) error {
	if len(batch) == 0 {
		return nil
	}

	// Build INSERT avec multiple VALUES
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
