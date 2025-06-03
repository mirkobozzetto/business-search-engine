package company

import (
	"csv-importer/api/cache"
	"database/sql"
	"log/slog"
	"os"
)

const MAX_COMPANIES = 100000

type companyService struct {
	db    *sql.DB
	cache *cache.RedisCache
}

func NewCompanyService(db *sql.DB) CompanyService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	redisCache := cache.NewRedisCache(cache.CacheConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	})

	return &companyService{
		db:    db,
		cache: redisCache,
	}
}
