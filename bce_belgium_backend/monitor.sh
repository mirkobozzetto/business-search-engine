#!/bin/bash

echo "📊 BCE System Status"
echo "===================="

echo "🐳 Docker Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "💾 Docker Volumes:"
docker volume ls --filter name=postgres

echo ""
echo "🔍 Database Connection:"
if docker exec bce_postgres pg_isready -U mirkobozzetto > /dev/null 2>&1; then
    echo "✅ PostgreSQL is running"
else
    echo "❌ PostgreSQL is down"
    exit 1
fi

echo ""
echo "📈 Table Statistics:"
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "
SELECT
    relname as \"Table\",
    n_live_tup as \"Rows\",
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||relname)) as \"Size\"
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;"

echo ""
echo "🌐 API Status:"
API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health)
if [ "$API_STATUS" = "200" ]; then
    echo "✅ API is running at http://localhost:8080"
else
    echo "❌ API is down (HTTP code: $API_STATUS)"
fi

echo ""
echo "💾 Disk Usage:"
docker system df

echo ""
echo "📊 Quick Data Summary:"
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "
SELECT
    count(*) as total_tables,
    sum(n_live_tup) as total_rows
FROM pg_stat_user_tables;"
