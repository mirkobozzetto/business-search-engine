package _cli

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func ListAvailableCSVs() {
	csvDir := "../bce_mai_2025"

	files, err := os.ReadDir(csvDir)
	if err != nil {
		log.Fatalf("❌ Cannot read directory %s: %v", csvDir, err)
	}

	fmt.Printf("📁 Available CSV files in %s:\n\n", csvDir)

	csvCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			tableName := GenerateTableName(file.Name())
			info, _ := file.Info()
			size := FormatFileSize(info.Size())
			fmt.Printf("   📄 %-20s → table '%s' (%s)\n", file.Name(), tableName, size)
			csvCount++
		}
	}

	if csvCount == 0 {
		fmt.Println("   No CSV files found")
	}
}
