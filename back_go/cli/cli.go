package cli

import (
	"csv-importer/cli/handlers"
	"database/sql"
	"fmt"
	"os"
)

type CLI struct {
	db *sql.DB
}

func New(db *sql.DB) *CLI {
	return &CLI{db: db}
}

func Run(db *sql.DB, args []string) {
	cli := New(db)
	cli.Execute(args)
}

func (c *CLI) Execute(args []string) {
	if len(args) < 2 {
		handlers.ShowHelp()
		os.Exit(1)
	}

	command := args[1]

	switch command {
	case "api":
		handlers.HandleAPI()
	case "all":
		handlers.HandleImportAll(c.db)
	case "list":
		handlers.HandleListCSVs()
	case "tables":
		handlers.HandleListTables(c.db)
	case "stats":
		handlers.HandleShowStats(c.db)
	case "info":
		handlers.HandleTableInfo(c.db, args[2:])
	case "columns":
		handlers.HandleShowColumns(c.db, args[2:])
	case "values":
		handlers.HandleColumnValues(c.db, args[2:])
	case "search":
		handlers.HandleSearch(c.db, args[2:])
	case "count":
		handlers.HandleCount(c.db, args[2:])
	case "sample":
		handlers.HandleSample(c.db, args[2:])
	case "export":
		handlers.HandleExport(c.db, args[2:])
	case "preview":
		handlers.HandlePreview(c.db, args[2:])
	case "help", "--help", "-h":
		handlers.ShowHelp()
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		handlers.ShowHelp()
		os.Exit(1)
	}
}
