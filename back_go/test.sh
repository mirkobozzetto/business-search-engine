#!/bin/bash
docker-compose down
docker volume prune -f
docker-compose up -d
docker logs bce_postgres
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "SELECT version();"
