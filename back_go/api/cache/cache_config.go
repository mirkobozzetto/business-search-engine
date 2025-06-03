package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

const (
	COMPRESSION_THRESHOLD = 1024              // 1KB
	MAX_UNCOMPRESSED_SIZE = 200 * 1024 * 1024 // 200MB pour données non-compressées
	MAX_COMPRESSED_SIZE   = 50 * 1024 * 1024  // 50MB pour données compressées
	COMPRESSION_PREFIX    = "gz:"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(config CacheConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
