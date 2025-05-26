package main

import (
	"csv-importer/cli"
	"csv-importer/config"
	"csv-importer/database"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("❌ DB connection failed:", err)
	}
	defer db.Close()

	fmt.Println("✅ DB connected")

	cli.Run(db, os.Args)
}
