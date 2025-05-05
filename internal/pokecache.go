package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	// Start the reap loop in background goroutine
	go c.reapLoop()

	return c
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	newCache := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	c.entries[key] = newCache
	c.mu.Unlock()
}

// Get retrieves an entry from the cached entries
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, found := c.entries[key]
	c.mu.Unlock()

	if !found {
		return nil, false
	}
	return entry.val, true
}

// Reaploop removes old entries
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	for {
		<-ticker.C
		c.mu.Lock()

		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.entries, key)
			}
		}

		c.mu.Unlock()
	}
}
