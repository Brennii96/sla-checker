package cache

import (
	"sync"
	"time"
)

// CacheEntry holds the data and its expiration time.
type CacheEntry[T any] struct {
	Data   T
	Expiry time.Time
}

// Cache is a generic cache structure that supports storing any type of data.
type Cache[T any] struct {
	data map[string]CacheEntry[T]
	mu   sync.Mutex
	ttl  time.Duration
}

// NewCache creates a new cache with a given TTL (time-to-live) duration.
func NewCache[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		data: make(map[string]CacheEntry[T]),
		ttl:  ttl,
	}
}

// Get retrieves an item from the cache. Returns the value and a boolean indicating if the value was found and not expired.
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.data[key]
	if !found || time.Now().After(entry.Expiry) {
		var zero T
		return zero, false
	}
	return entry.Data, true
}

// Set stores an item in the cache with the specified key and sets its expiration time based on the TTL.
func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheEntry[T]{
		Data:   value,
		Expiry: time.Now().Add(c.ttl),
	}
}

// Delete removes an item from the cache.
func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
