package cache

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

func (r *RedisCache) Set(key string, value any, ttl time.Duration) error {
	start := time.Now()
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	result, err := r.smartCompress(key, jsonData)
	if err != nil {
		return err
	}
	err = r.client.Set(r.ctx, result.Key, result.Data, ttl).Err()
	duration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to cache data: %v", err)
	}
	slog.Info("Cache write",
		"key", key,
		"original_mb", result.OriginalSize/1024/1024,
		"final_mb", result.FinalSize/1024/1024,
		"duration_ms", duration.Milliseconds(),
		"compressed", result.IsCompressed)
	return nil
}

func (r *RedisCache) Get(key string, dest any) error {
	compressedKey := COMPRESSION_PREFIX + key
	val, err := r.client.Get(r.ctx, compressedKey).Result()
	isCompressed := true
	if err == redis.Nil {
		val, err = r.client.Get(r.ctx, key).Result()
		isCompressed = false
	}
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get key: %v", err)
	}
	var jsonData []byte
	if isCompressed {
		decompressed, err := r.decompressData([]byte(val))
		if err != nil {
			return fmt.Errorf("failed to decompress data: %v", err)
		}
		jsonData = decompressed
	} else {
		jsonData = []byte(val)
	}
	return json.Unmarshal(jsonData, dest)
}
