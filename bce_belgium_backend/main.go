package main

import (
	"csv-importer/cli"
	"csv-importer/config"
	"csv-importer/database"
	"log/slog"
	"os"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("❌ DB connection failed", "error", err)
	}
	defer db.Close()

	slog.Info("✅ DB connected")

	cli.Run(db, os.Args)
}
