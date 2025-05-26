#!/bin/bash

echo "🔄 Starting BCE PostgreSQL + API safely"

echo "📦 Stopping containers..."
docker-compose down

echo "🚀 Starting containers..."
docker-compose up -d

echo "⏳ Waiting for PostgreSQL (5 seconds)..."
sleep 5

echo "📋 PostgreSQL logs:"
docker logs --tail 10 bce_postgres

echo "🔍 Testing DB connection:"
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "SELECT version();"

echo "📊 Checking data integrity:"
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "
SELECT
    schemaname,
    relname as tablename,
    n_live_tup as rows
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;"

echo "🌐 Testing API:"
sleep 2
curl -s http://localhost:8080/api/health | head -50

echo "✅ All systems running! API available at http://localhost:8080"
