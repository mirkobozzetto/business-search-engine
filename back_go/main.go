package main

import (
	"csv-importer/config"
	"csv-importer/csv"
	"csv-importer/database"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	fmt.Println("✅ DB connected")

	csvPath := "../bce_mai_2025/activity.csv"
	tableName := "activity"

	if err := csv.ProcessCSV(db, csvPath, tableName); err != nil {
		log.Fatal("CSV processing failed:", err)
	}

	fmt.Println("🎉 Import terminé!")
}
