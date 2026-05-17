package main

import (
	"log/slog"
	"os"
	"sirene-importer/cli"
	"sirene-importer/config"
	"sirene-importer/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("DB connection failed", "error", err)
	}
	defer func() { _ = db.Close() }()

	slog.Info("DB connected")

	cli.Run(db, os.Args)
}
