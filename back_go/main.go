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
	fmt.Printf("📊 Colonnes: %v\n", headers)

	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	if _, err := db.Exec(dropSQL); err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	var columns []string
	for _, header := range headers {
		cleanHeader := strings.ReplaceAll(header, " ", "_")
		cleanHeader = strings.ReplaceAll(cleanHeader, "-", "_")
		cleanHeader = strings.ToLower(cleanHeader)
		columns = append(columns, cleanHeader+" VARCHAR(255)")
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))
	fmt.Printf("🏗️  Création table: %s\n", tableName)

	if _, err := db.Exec(createSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	placeholders := make([]string, len(headers))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s VALUES (%s)",
		tableName, strings.Join(placeholders, ", "))

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("error preparing INSERT: %v", err)
	}
	defer stmt.Close()

	lineCount := 0
	batchSize := 1000

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading line %d: %v", lineCount+1, err)
		}

		// Covert []string en any (or []interface{})
		values := make([]any, len(record))
		for i, v := range record {
			values[i] = v
		}

		if _, err := stmt.Exec(values...); err != nil {
			return fmt.Errorf("error inserting line %d: %v", lineCount+1, err)
		}

		lineCount++
		if lineCount%batchSize == 0 {
			fmt.Printf("📈 Processed: %d lines\n", lineCount)
		}
	}

	fmt.Printf("✅ Total inserted: %d lines\n", lineCount)
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
