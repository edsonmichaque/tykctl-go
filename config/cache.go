package config

import (
	"fmt"
	"sync"
	"time"
)

// CacheOptions provides configuration for cache
type CacheOptions struct {
	TTL         time.Duration `json:"ttl"`
	MaxSize     int64         `json:"max_size"`
	Compression bool          `json:"compression"`
	Encryption  bool          `json:"encryption"`
}

// NewCache creates a new cache instance
func NewCache(opts CacheOptions) (Cache, error) {
	return &memoryCache{
		data: make(map[string]cacheItem),
		ttl:  opts.TTL,
		mu:   sync.RWMutex{},
	}, nil
}

// memoryCache implements a simple in-memory cache
type memoryCache struct {
	data map[string]cacheItem
	ttl  time.Duration
	mu   sync.RWMutex
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func (c *memoryCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	if time.Now().After(item.expiresAt) {
		delete(c.data, key)
		return nil, fmt.Errorf("key expired")
	}

	return item.value, nil
}

func (c *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl == 0 {
		ttl = c.ttl
	}

	c.data[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (c *memoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}

func (c *memoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]cacheItem)
	return nil
}

func (c *memoryCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = nil
	return nil
}

func (c *memoryCache) TTL() time.Duration {
	return c.ttl
}