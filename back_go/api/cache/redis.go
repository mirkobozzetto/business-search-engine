package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
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

func (r *RedisCache) Set(key string, value any, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	return r.client.Set(r.ctx, key, jsonData, ttl).Err()
}

func (r *RedisCache) Get(key string, dest any) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get key: %v", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

func (r *RedisCache) GetKeys(pattern string) ([]string, error) {
	return r.client.Keys(r.ctx, pattern).Result()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
