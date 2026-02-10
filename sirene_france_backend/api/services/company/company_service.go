package company

import (
	"database/sql"
	"os"
	"sirene-importer/api/cache"
)

type companyService struct {
	db    *sql.DB
	cache *cache.RedisCache
}

func NewCompanyService(db *sql.DB) *companyService {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6380"
	}
	redisCache := cache.NewRedisCache(cache.CacheConfig{
		Host: redisHost,
		Port: redisPort,
		DB:   0,
	})
	return &companyService{db: db, cache: redisCache}
}
