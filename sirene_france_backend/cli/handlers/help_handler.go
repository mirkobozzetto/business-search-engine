package handlers

import "fmt"

func ShowHelp() {
	fmt.Println(`
SIRENE France Backend

Usage: sirene-api <command>

Commands:
  api                    Start the API server (port 8081)
  all                    Import all ZIP files from ../sirene_data
  tables                 List all database tables

  help                   Show this help message

API Endpoints:
  GET /api/health
  GET /api/companies/search/naf?code={code}&limit={limit}
  GET /api/companies/search/denomination?q={query}&limit={limit}
  GET /api/companies/search/codepostal?q={code}&limit={limit}
  GET /api/companies/search/commune?q={commune}&limit={limit}
  GET /api/companies/search/etatadministratif?q={etat}&limit={limit}
  GET /api/companies/search/datecreation?from={date}&to={date}&limit={limit}
  GET /api/companies/search/multi?naf={code}&denomination={query}&codepostal={code}&commune={commune}&etat={etat}&from={date}&to={date}&limit={limit}`)
}
