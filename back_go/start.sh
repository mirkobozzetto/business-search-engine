#!/bin/bash

echo "🔄 Starting BCE PostgreSQL + Redis + API safely"

echo "🔍 Checking containers..."
if ! docker ps | grep -q bce_postgres; then
    echo "🚀 Starting containers..."
    docker-compose up -d
else
    echo "✅ Containers already running"
fi

echo "⏳ Waiting for PostgreSQL to be ready..."
until docker exec bce_postgres pg_isready -U mirkobozzetto -d bce_db > /dev/null 2>&1; do
    sleep 2
    echo "   Still waiting..."
done
echo "✅ PostgreSQL ready!"

echo "📋 PostgreSQL logs:"
docker logs --tail 10 bce_postgres

echo "📋 Redis logs:"
docker logs --tail 10 bce_redis

echo "🔍 Testing DB connection:"
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "SELECT version();"

echo "🔍 Testing Redis connection:"
docker exec bce_redis redis-cli ping

echo "📊 Checking data integrity:"
docker exec bce_postgres psql -U mirkobozzetto -d bce_db -c "
SELECT
    schemaname,
    relname as tablename,
    n_live_tup as rows
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;"

echo "📊 Redis info:"
docker exec bce_redis redis-cli info server | grep redis_version

echo "🌐 Testing API:"
sleep 2
curl -s http://localhost:8080/api/health | head -50

echo "✅ All systems running!"
echo "   PostgreSQL: http://localhost:5433"
echo "   Redis: http://localhost:6379"
echo "   API: http://localhost:8080"
