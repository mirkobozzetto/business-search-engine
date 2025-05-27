#!/bin/bash

echo "🔄 Starting BCE PostgreSQL + Redis + API safely"

echo "📦 Stopping containers..."
docker-compose down

echo "🚀 Starting containers..."
docker-compose up -d

echo "⏳ Waiting for services (7 seconds)..."
sleep 7

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
