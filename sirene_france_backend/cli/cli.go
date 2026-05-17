package cli

import (
	"database/sql"
	"os"
	"sirene-importer/cli/handlers"
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
	case "tables":
		handlers.HandleListTables(c.db)
	case "indexes":
		handlers.HandleCreateIndexes(c.db)
	case "naf":
		handlers.HandleImportNaf(c.db)
	case "help", "--help", "-h":
		handlers.ShowHelp()
	default:
		handlers.ShowHelp()
		os.Exit(1)
	}
}
