package cache

import (
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RedisCache) Delete(key string) error {
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
	keys1, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	keys2, err := r.client.Keys(r.ctx, COMPRESSION_PREFIX+pattern).Result()
	if err != nil {
		return nil, err
	}

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
