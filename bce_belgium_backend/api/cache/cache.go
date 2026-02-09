package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data    any
	Expires time.Time
}

type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.data[key]
	if !exists || time.Now().After(item.Expires) {
			return nil, false
	}
	return item.Data, true
}
