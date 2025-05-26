# PostgreSQL Docker Setup

## Quick Start

```bash
docker-compose up -d
sleep 15
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "SELECT version();"
```

## Port Config

```yaml
ports:
  - "5433:5432" # Use 5433 if local postgres on 5432
  - "5432:5432" # Use 5432 if no local postgres
```

## Commands

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Reset (deletes data)
docker-compose down
docker volume prune -f
docker-compose up -d

# Connect
docker exec -it bce_postgres psql -U username -d database

# Logs
docker logs bce_postgres
```

## Troubleshooting

**Env vars not working?** → Check `.env` file location and `docker-compose config`

**Port conflict?** → Change host port in docker-compose.yml

**Old data?** → Reset with volume prune

## Client Connections

### Using psql (inside container)

```bash
docker exec -it bce_postgres psql -U mirkobozzetto -d bce_db
```

### Using psql (from host)

```bash
# Port 5433
psql -h localhost -p 5433 -U your_username -d your_database_name

# Port 5432
psql -h localhost -p 5432 -U your_username -d your_database_name
```

### Using pgcli (from host)

```bash
# Port 5433
pgcli -h localhost -p 5433 -U your_username -d your_database_name

# Port 5432
pgcli -h localhost -p 5432 -U your_username -d your_database_name
```

### Connection String Format

```
postgresql://username:password@localhost:5433/database_name
```

## Management Commands

### Start services

```bash
docker-compose up -d
```

### Stop services

```bash
docker-compose down
```

### View logs

```bash
docker logs bce_postgres
```

### Reset database (WARNING: Deletes all data)

```bash
docker-compose down
docker volume rm your_project_postgres_data
docker-compose up -d
```

## Security Notes

- Never use weak passwords in production
- Don't commit `.env` files to version control
- Consider using Docker secrets for production deployments
- Regularly update the PostgreSQL image

## Example .env for Development

```bash
POSTGRES_DB=development_db
POSTGRES_USER=dev_user
POSTGRES_PASSWORD=dev_password_123
```

## Example .env for Production

```bash
POSTGRES_DB=production_db
POSTGRES_USER=app_user
POSTGRES_PASSWORD=very_secure_random_password_here
```
