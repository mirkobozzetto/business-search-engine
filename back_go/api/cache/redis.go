package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
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

const (
	COMPRESSION_THRESHOLD = 1024             // 1KB
	MAX_CACHE_SIZE        = 50 * 1024 * 1024 // 50MB
	COMPRESSION_PREFIX    = "gz:"
)

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
	start := time.Now()

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	originalSize := len(jsonData)
	var finalKey string
	var finalData []byte

	if originalSize > MAX_CACHE_SIZE {
		slog.Warn("Dataset too large for cache, skipping",
			"key", key,
			"size_mb", originalSize/1024/1024,
			"max_mb", MAX_CACHE_SIZE/1024/1024)
		return fmt.Errorf("dataset too large: %d MB", originalSize/1024/1024)
	}

	if originalSize > COMPRESSION_THRESHOLD {
		compressed, err := r.compressData(jsonData)
		if err != nil {
			slog.Warn("Compression failed, using uncompressed",
				"key", key,
				"error", err)
			finalKey = key
			finalData = jsonData
		} else {
			compressionRatio := float64(len(compressed)) / float64(originalSize) * 100
			finalKey = COMPRESSION_PREFIX + key
			finalData = compressed

			slog.Info("Data compressed for cache",
				"key", key,
				"original_kb", originalSize/1024,
				"compressed_kb", len(compressed)/1024,
				"compression_ratio", fmt.Sprintf("%.1f%%", compressionRatio))
		}
	} else {
		finalKey = key
		finalData = jsonData
	}

	err = r.client.Set(r.ctx, finalKey, finalData, ttl).Err()
	duration := time.Since(start)

	if err != nil {
		slog.Error("Cache write failed",
			"key", key,
			"size_kb", originalSize/1024,
			"duration_ms", duration.Milliseconds(),
			"compressed", strings.HasPrefix(finalKey, COMPRESSION_PREFIX),
			"error", err.Error())
		return fmt.Errorf("failed to cache data: %v", err)
	}

	slog.Info("Cache write successful",
		"key", key,
		"size_kb", originalSize/1024,
		"final_size_kb", len(finalData)/1024,
		"duration_ms", duration.Milliseconds(),
		"compressed", strings.HasPrefix(finalKey, COMPRESSION_PREFIX))

	return nil
}

func (r *RedisCache) Get(key string, dest any) error {
	start := time.Now()

	// Try compressed version first
	compressedKey := COMPRESSION_PREFIX + key
	val, err := r.client.Get(r.ctx, compressedKey).Result()
	isCompressed := true

	if err == redis.Nil {
		// Try uncompressed version
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
		"size_kb", len(jsonData)/1024,
		"compressed", isCompressed)

	return nil
}

func (r *RedisCache) Delete(key string) error {
	// Delete both compressed and uncompressed versions
	pipe := r.client.Pipeline()
	pipe.Del(r.ctx, key)
	pipe.Del(r.ctx, COMPRESSION_PREFIX+key)

	results, err := pipe.Exec(r.ctx)
	if err != nil {
		return err
	}

	deletedCount := 0
	for _, result := range results {
		if result.Err() == nil {
			if delResult, ok := result.(*redis.IntCmd); ok {
				deletedCount += int(delResult.Val())
			}
		}
	}

	slog.Info("Cache delete",
		"key", key,
		"deleted_keys", deletedCount)

	return nil
}

func (r *RedisCache) Exists(key string) (bool, error) {
	// Check both compressed and uncompressed versions
	pipe := r.client.Pipeline()
	cmd1 := pipe.Exists(r.ctx, key)
	cmd2 := pipe.Exists(r.ctx, COMPRESSION_PREFIX+key)

	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return false, err
	}

	return cmd1.Val() > 0 || cmd2.Val() > 0, nil
}

func (r *RedisCache) GetKeys(pattern string) ([]string, error) {
	// Get both compressed and uncompressed keys
	keys1, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	keys2, err := r.client.Keys(r.ctx, COMPRESSION_PREFIX+pattern).Result()
	if err != nil {
		return nil, err
	}

	// Remove compression prefix and deduplicate
	keySet := make(map[string]bool)
	for _, key := range keys1 {
		keySet[key] = true
	}

	for _, key := range keys2 {
		cleanKey := strings.TrimPrefix(key, COMPRESSION_PREFIX)
		keySet[cleanKey] = true
	}

	result := make([]string, 0, len(keySet))
	for key := range keySet {
		result = append(result, key)
	}

	return result, nil
}

func (r *RedisCache) GetCacheStats() map[string]any {
	info := r.client.Info(r.ctx, "memory")

	return map[string]any{
		"redis_info": info.Val(),
		"timestamp":  time.Now(),
	}
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *RedisCache) decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
