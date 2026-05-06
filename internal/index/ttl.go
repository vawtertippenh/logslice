package index

import (
	"sync"
	"time"
)

// TTLEntry holds a cached value along with its expiry time.
type TTLEntry struct {
	value     interface{}
	expiresAt time.Time
}

// IsExpired reports whether the entry has passed its expiry time.
func (e TTLEntry) IsExpired(now time.Time) bool {
	return now.After(e.expiresAt)
}

// TTLCache is a thread-safe cache that evicts entries after a configurable
// time-to-live duration. It is intended to wrap the index Cache to prevent
// stale index data from persisting across log rotations.
type TTLCache struct {
	mu      sync.Mutex
	entries map[string]TTLEntry
	ttl     time.Duration
}

// NewTTLCache creates a TTLCache with the given time-to-live per entry.
func NewTTLCache(ttl time.Duration) *TTLCache {
	return &TTLCache{
		entries: make(map[string]TTLEntry),
		ttl:     ttl,
	}
}

// Put stores a value under key with an expiry of now+ttl.
func (c *TTLCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = TTLEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value by key. The second return value is false when the
// key is absent or its entry has expired (in which case the entry is deleted).
func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if e.IsExpired(time.Now()) {
		delete(c.entries, key)
		return nil, false
	}
	return e.value, true
}

// Delete removes an entry from the cache.
func (c *TTLCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Purge removes all expired entries and returns the number of entries removed.
func (c *TTLCache) Purge() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	removed := 0
	for k, e := range c.entries {
		if e.IsExpired(now) {
			delete(c.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of entries currently in the cache (including expired).
func (c *TTLCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
