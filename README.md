# Business Search Engine

Multi-country business search engine for **Belgian** (BCE) and **French** (SIRENE) company registries.
120M+ rows across two Go APIs with PostgreSQL, Redis caching, and a Next.js frontend.

## Architecture

```
business_search_engine/
├── bce_belgium_backend/        Go API — Belgian BCE registry (47M rows)
├── sirene_france_backend/      Go API — French SIRENE registry (72M rows)
└── sirene_france_frontend/     Next.js 16 — Search UI for French companies
```

```
┌──────────────────────┐         ┌──────────────────────┐
│  Belgium (BCE)       │         │  France (SIRENE)     │
│  Go API :8080        │         │  Go API :8081        │
│  PostgreSQL :5433    │         │  PostgreSQL :5434    │
│  Redis :6379         │         │  Redis :6380         │
│  47M rows / 10 tables│         │  72M rows / 3 tables │
└──────────────────────┘         └──────────────────────┘
                                          ▲
                                          │
                                 ┌────────┴────────┐
                                 │  Next.js :3000   │
                                 │  React 19 + TS   │
                                 └─────────────────┘
```

## Tech Stack

| Layer    | Technology                                      |
|----------|------------------------------------------------|
| Backend  | Go 1.24, Gin, pgx                              |
| Database | PostgreSQL 15 (pg_trgm, unaccent)              |
| Cache    | Redis 7 (gzip compression, 24h TTL, AOF)       |
| Frontend | Next.js 16, React 19, TypeScript, Tailwind CSS 4 |
| UI       | shadcn/ui, TanStack React Query 5              |
| Infra    | Docker Compose                                 |

## Features

### Search Capabilities

| Feature             | Belgium (BCE)   | France (SIRENE)          |
|---------------------|-----------------|--------------------------|
| Activity code       | NACE (5 digits) | NAF (XX.XXZ)             |
| Company name        | Multi-word AND  | Multi-word AND + trigram  |
| Postal code         | Yes             | Yes                      |
| City / commune      | —               | Yes (fuzzy)              |
| Administrative state| —               | Active / Ceased          |
| Creation date range | Yes             | Yes                      |
| Multi-criteria      | Yes             | Yes                      |
| Direct lookup       | Entity number   | SIREN / SIRET            |
| Contact info        | Email, phone, web, fax | —                  |
| CSV / JSON export   | Yes             | JSON                     |
| Accent-insensitive  | —               | Yes (immutable_unaccent) |

### Performance

- **French import**: 42.7M rows in ~10 min (~50k rows/sec, 8 workers)
- **Redis caching**: gzip-compressed, 24h TTL, 200 MB decompression limit
- **PostgreSQL**: trigram indexes, custom `immutable_unaccent()`, tuned for 72M rows

## Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Node.js 20+ and pnpm (for frontend)

## Quick Start

### Belgium (BCE)

```bash
# Download BCE CSV files from https://kbopub.economie.fgov.be/kbo-open-data
# Place them in bce_belgium_backend/

make bce-up          # Start PostgreSQL + Redis
make bce-import      # Import CSV files (~47M rows)
make bce-api         # Start API on :8080
```

### France (SIRENE)

```bash
# Download SIRENE ZIP files from https://www.data.gouv.fr/fr/datasets/base-sirene-des-entreprises-et-de-leurs-etablissements-siren-siret/
# Place them in sirene_data/

make sirene-up       # Start PostgreSQL + Redis
make sirene-import   # Import ZIP files (~72M rows, ~10 min)
make sirene-api      # Start API on :8081
make sirene-front-dev  # Start frontend on :3000
```

### Both

```bash
make up-all          # Start all containers
make down-all        # Stop everything
make ps              # Show container status
make help            # All available commands
```

## API Endpoints

### Belgium — `:8080`

```bash
GET /api/companies/search/nace?code=62010
GET /api/companies/search/denomination?q=informatique
GET /api/companies/search/zipcode?q=1000
GET /api/companies/search/startdate?from=01-01-2025
GET /api/companies/search/multi?nace=62010&zipcode=1000

GET /api/tables
GET /api/data/:table/preview
GET /api/export/:table
GET /api/health
```

### France — `:8081`

```bash
GET /api/companies/lookup/:identifier          # SIREN or SIRET
GET /api/companies/search/naf?code=62.01Z
GET /api/companies/search/denomination?q=google
GET /api/companies/search/codepostal?q=75001
GET /api/companies/search/commune?q=paris
GET /api/companies/search/etatadministratif?q=A
GET /api/companies/search/datecreation?from=2025-01-01&to=2025-12-31
GET /api/companies/search/multi?naf=62.01Z&commune=paris&etat=A

GET /api/naf/search?q=informatique
GET /api/naf/sections
GET /api/naf/code/:code
GET /api/health
```

## Data Sources

| Country | Source | Rows | Tables | Update Frequency |
|---------|--------|------|--------|-----------------|
| Belgium | [BCE Open Data](https://kbopub.economie.fgov.be/kbo-open-data) | 47M | 10 | Monthly |
| France  | [INSEE SIRENE](https://www.data.gouv.fr/fr/datasets/base-sirene-des-entreprises-et-de-leurs-etablissements-siren-siret/) | 72M | 3 | Monthly |

## Code Quality

```bash
make format    # gofmt on all Go code
make lint      # golangci-lint on all Go code
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[Apache License 2.0](LICENCE) — Copyright 2025 Mirko Bozzetto
