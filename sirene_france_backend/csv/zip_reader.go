package csv

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type ZIPFile struct {
	Name      string
	Path      string
	TableName string
}

func getZIPFiles(zipDir string) ([]ZIPFile, error) {
	var zipFiles []ZIPFile

	entries, err := filepath.Glob(filepath.Join(zipDir, "*.zip"))
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		base := filepath.Base(entry)
		tableName := deriveTableName(base)
		zipFiles = append(zipFiles, ZIPFile{
			Name:      base,
			Path:      entry,
			TableName: tableName,
		})
	}

	return zipFiles, nil
}

func deriveTableName(zipName string) string {
	name := strings.TrimSuffix(zipName, ".zip")
	name = strings.TrimSuffix(name, "_utf8")

	switch {
	case strings.Contains(strings.ToLower(name), "unitelegale"):
		return "unite_legale"
	case strings.Contains(strings.ToLower(name), "etablissement"):
		return "etablissement"
	default:
		name = strings.ToLower(name)
		name = strings.ReplaceAll(name, "-", "_")
		name = strings.ReplaceAll(name, " ", "_")
		return name
	}
}

func ProcessAllZIPs(db *sql.DB, zipDir string) error {
	zipFiles, err := getZIPFiles(zipDir)
	if err != nil {
		return fmt.Errorf("error scanning ZIP directory: %v", err)
	}

	if len(zipFiles) == 0 {
		return fmt.Errorf("no ZIP files found in %s", zipDir)
	}

	fmt.Printf("Processing %d ZIP files\n", len(zipFiles))

	totalStart := time.Now()

	for i, zf := range zipFiles {
		fmt.Printf("\n[%d/%d] Processing %s...\n", i+1, len(zipFiles), zf.Name)

		if err := ProcessZIPFile(db, zf.Path, zf.TableName); err != nil {
			return fmt.Errorf("error processing %s: %v", zf.Name, err)
		}

		fmt.Printf("[%d/%d] %s -> table '%s' completed\n", i+1, len(zipFiles), zf.Name, zf.TableName)
	}

	totalElapsed := time.Since(totalStart)
	fmt.Printf("\nAll %d tables created in %.2f minutes\n", len(zipFiles), totalElapsed.Minutes())

	fmt.Println("\nCr\u00e9ation des indexes...")
	if err := CreateIndexes(db); err != nil {
		fmt.Printf("Erreur indexes: %v\n", err)
	}

	return nil
}

func ProcessZIPFile(db *sql.DB, zipPath, tableName string) error {
	start := time.Now()

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip %s: %w", zipPath, err)
	}
	defer r.Close()

	for _, f := range r.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".csv") {
			continue
		}

		fmt.Printf("Found CSV in ZIP: %s (%.2f GB)\n", f.Name, float64(f.UncompressedSize64)/(1024*1024*1024))

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open CSV in zip: %w", err)
		}
		defer rc.Close()

		csvReader := csv.NewReader(rc)
		csvReader.Comma = ','
		csvReader.LazyQuotes = true

		headers, err := csvReader.Read()
		if err != nil {
			return fmt.Errorf("error reading headers: %w", err)
		}

		fmt.Printf("CSV: %s (%d columns)\n", f.Name, len(headers))

		cleanHeaders, columns := PrepareHeaders(headers)

		if err := setupTable(db, tableName, columns); err != nil {
			return err
		}

		rc.Close()

		rc2, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to reopen CSV: %w", err)
		}
		defer rc2.Close()

		totalLines, err := ProcessPipelineFromReader(rc2, tableName, cleanHeaders)
		if err != nil {
			return err
		}

		elapsed := time.Since(start)
		linesPerSec := float64(totalLines) / elapsed.Seconds()

		fmt.Printf("Done: %d lines in %.2f sec (%.0f lines/sec)\n",
			totalLines, elapsed.Seconds(), linesPerSec)

		return nil
	}

	return fmt.Errorf("no CSV file found in %s", zipPath)
}

func setupTable(db *sql.DB, tableName string, columns []string) error {
	OptimizeForBulkInsert(db)

	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	if _, err := db.Exec(dropSQL); err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))
	fmt.Printf("Creating table: %s\n", tableName)

	if _, err := db.Exec(createSQL); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}
