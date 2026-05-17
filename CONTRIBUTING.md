# Contributing to Business Search Engine

This project is open source under the [Apache License 2.0](LICENCE).
Contributions are welcome.

## Before you start

- Open an issue to discuss the feature or bug before writing code
- One PR per feature or fix — keep it focused

## Setup

```bash
# Go backends
go install golang.org/dl/go1.24@latest
cp bce_belgium_backend/.env.example bce_belgium_backend/.env
cp sirene_france_backend/.env.example sirene_france_backend/.env

# Frontend
cd sirene_france_frontend && pnpm install

# Verify
make format
make lint
```

## Code style

### Go

- `gofmt` before every commit
- `golangci-lint run ./...` must pass with no warnings
- No comments in code — use clear naming instead
- Follow existing patterns in the codebase

### Frontend (TypeScript / React)

- Follow existing component structure in `src/components/`
- One component per file
- Use existing hooks in `src/hooks/` as reference

## Pull requests

1. Fork the repo and create your branch from `main`
2. Run `make format` and `make lint`
3. Test your changes locally against real data
4. Open a PR against `main`

## Project structure

```
bce_belgium_backend/     Go API for Belgian BCE registry
sirene_france_backend/   Go API for French SIRENE registry
sirene_france_frontend/  Next.js frontend for French search
```

Each backend follows the same pattern:
- `api/services/company/` — search logic
- `api/models/` — data structures
- `api/cache/` — Redis operations
- `csv/` — import pipeline

## Database

Do not modify database schemas without discussion in an issue first.
Both backends use PostgreSQL with specialized indexes (trigram, unaccent).

## License

By contributing, you agree that your contributions will be licensed under the [Apache License 2.0](LICENCE).
