package csv

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CSVFile struct {
	Name      string
	Path      string
	TableName string
}

func ProcessAllCSVs(db *sql.DB, csvDir string) error {
	csvFiles, err := getCSVFiles(csvDir)
	if err != nil {
		return fmt.Errorf("error scanning CSV directory: %v", err)
	}

	if len(csvFiles) == 0 {
		return fmt.Errorf("no CSV files found in %s", csvDir)
	}

	fmt.Printf("📁 Found %d CSV files to process\n", len(csvFiles))

	totalStart := time.Now()

	for i, csvFile := range csvFiles {
		fmt.Printf("\n🔄 [%d/%d] Processing %s...\n", i+1, len(csvFiles), csvFile.Name)

		if err := ProcessCSV(db, csvFile.Path, csvFile.TableName); err != nil {
			return fmt.Errorf("error processing %s: %v", csvFile.Name, err)
		}

		fmt.Printf("✅ [%d/%d] %s → table '%s' completed\n", i+1, len(csvFiles), csvFile.Name, csvFile.TableName)
	}

	totalElapsed := time.Since(totalStart)
	fmt.Printf("\n🏆 ALL %d TABLES CREATED in %.2f minutes\n", len(csvFiles), totalElapsed.Minutes())

	// Show summary
	fmt.Println("\n📊 SUMMARY:")
	for _, csvFile := range csvFiles {
		fmt.Printf("   • %s → table '%s'\n", csvFile.Name, csvFile.TableName)
	}

	return nil
}

func getCSVFiles(csvDir string) ([]CSVFile, error) {
	var csvFiles []CSVFile

	files, err := os.ReadDir(csvDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			continue
		}

		// Generate table name from filename
		tableName := strings.TrimSuffix(strings.ToLower(file.Name()), ".csv")
		tableName = strings.ReplaceAll(tableName, "-", "_")
		tableName = strings.ReplaceAll(tableName, " ", "_")

		csvFiles = append(csvFiles, CSVFile{
			Name:      file.Name(),
			Path:      filepath.Join(csvDir, file.Name()),
			TableName: tableName,
		})
	}

	return csvFiles, nil
}
