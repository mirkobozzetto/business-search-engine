package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RedisCache) Set(key string, value any, ttl time.Duration) error {
	start := time.Now()

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	result, err := r.smartCompress(key, jsonData)
	if err != nil {
		slog.Error("Smart compression failed",
			"key", key,
			"original_mb", len(jsonData)/1024/1024,
			"error", err)
		return err
	}

	err = r.client.Set(r.ctx, result.Key, result.Data, ttl).Err()
	duration := time.Since(start)

	if err != nil {
		slog.Error("Cache write failed",
			"key", key,
			"size_mb", result.OriginalSize/1024/1024,
			"final_size_mb", result.FinalSize/1024/1024,
			"duration_ms", duration.Milliseconds(),
			"compressed", result.IsCompressed,
			"error", err.Error())
		return fmt.Errorf("failed to cache data: %v", err)
	}

	slog.Info("Cache write successful",
		"key", key,
		"original_mb", result.OriginalSize/1024/1024,
		"final_size_mb", result.FinalSize/1024/1024,
		"duration_ms", duration.Milliseconds(),
		"compressed", result.IsCompressed,
		"compression_ratio", fmt.Sprintf("%.1f%%", result.CompressionRatio))

	return nil
}

func (r *RedisCache) Get(key string, dest any) error {
	start := time.Now()

	compressedKey := COMPRESSION_PREFIX + key
	val, err := r.client.Get(r.ctx, compressedKey).Result()
	isCompressed := true

	if err == redis.Nil {
		val, err = r.client.Get(r.ctx, key).Result()
		isCompressed = false
	}

	if err != nil {
		duration := time.Since(start)
		slog.Info("Cache miss",
			"key", key,
			"duration_ms", duration.Milliseconds())

		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get key: %v", err)
	}

	var jsonData []byte
	if isCompressed {
		decompressed, err := r.decompressData([]byte(val))
		if err != nil {
			slog.Error("Decompression failed",
				"key", key,
				"error", err)
			return fmt.Errorf("failed to decompress data: %v", err)
		}
		jsonData = decompressed
	} else {
		jsonData = []byte(val)
	}

	err = json.Unmarshal(jsonData, dest)
	duration := time.Since(start)

	if err != nil {
		slog.Error("Cache unmarshal failed",
			"key", key,
			"duration_ms", duration.Milliseconds(),
			"compressed", isCompressed,
			"error", err)
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	slog.Info("Cache hit",
		"key", key,
		"duration_ms", duration.Milliseconds(),
		"size_mb", len(jsonData)/1024/1024,
		"compressed", isCompressed)

	return nil
}
