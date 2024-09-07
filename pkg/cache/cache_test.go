package cache_test

import (
	"testing"
	"time"

	"github.com/brennii96/sla-checker/pkg/cache"
)

func TestCacheSetAndGet(t *testing.T) {
	c := cache.NewCache[string](5 * time.Second)

	// Test setting and getting a value
	c.Set("key1", "value1")

	value, found := c.Get("key1")
	if !found || value != "value1" {
		t.Errorf("expected value1, got %v", value)
	}

	// Test that the value is still there after some time but before expiration
	time.Sleep(2 * time.Second)
	value, found = c.Get("key1")
	if !found || value != "value1" {
		t.Errorf("expected value1, got %v", value)
	}
}

func TestCacheExpiration(t *testing.T) {
	c := cache.NewCache[string](2 * time.Second)

	// Set a value and let it expire
	c.Set("key1", "value1")

	time.Sleep(3 * time.Second) // wait for the cache to expire

	_, found := c.Get("key1")
	if found {
		t.Error("expected key1 to be expired and not found")
	}
}

func TestCacheDelete(t *testing.T) {
	c := cache.NewCache[string](5 * time.Second)

	// Test deleting a value
	c.Set("key1", "value1")
	c.Delete("key1")

	_, found := c.Get("key1")
	if found {
		t.Error("expected key1 to be deleted and not found")
	}
}
