package handlers

import "fmt"

func ShowHelp() {
	fmt.Println(`
SIRENE France Backend - API de recherche d'entreprises françaises

Usage: sirene-api <commande>

Commandes:
  api                    Démarrer le serveur API (port 8081)
  all                    Importer tous les fichiers ZIP depuis ../sirene_data
  indexes                Créer les indexes PostgreSQL (btree + trigram)
  tables                 Lister les tables de la base de données
  help                   Afficher cette aide

Endpoints API (port 8081):
  GET /api/health
  GET /api/companies/search/naf?code={code}&limit={n}&offset={n}
  GET /api/companies/search/denomination?q={query}&limit={n}&offset={n}
  GET /api/companies/search/codepostal?q={code}&limit={n}&offset={n}
  GET /api/companies/search/commune?q={commune}&limit={n}&offset={n}
  GET /api/companies/search/etatadministratif?q={A|C}&limit={n}&offset={n}
  GET /api/companies/search/datecreation?from={YYYY-MM-DD}&to={YYYY-MM-DD}&limit={n}&offset={n}
  GET /api/companies/search/multi?naf={code}&denomination={q}&codepostal={cp}&commune={c}&etat={A|C}&from={date}&to={date}&limit={n}&offset={n}

Paramètres de pagination:
  limit                  Nombre de résultats par page (défaut: 100, max: 10000)
  offset                 Position de départ dans les résultats (défaut: 0)

Exemples:
  curl "localhost:8081/api/companies/search/naf?code=62.01Z&limit=10"
  curl "localhost:8081/api/companies/search/denomination?q=google&limit=20&offset=40"
  curl "localhost:8081/api/companies/search/commune?q=paris&limit=50"
  curl "localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&etat=A&limit=10"`)
}
