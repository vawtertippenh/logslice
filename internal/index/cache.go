package index

import (
	"os"
	"sync"
	"time"
)

// entry holds a cached Index along with its file modification time.
type entry struct {
	index   *Index
	modTime time.Time
}

// Cache stores built indexes keyed by file path and evicts entries using LRU
// when the cache exceeds its configured capacity.
type Cache struct {
	mu      sync.Mutex
	entries map[string]*entry
	evict   *LRUEviction
}

// NewCache creates a Cache that holds at most capacity indexes in memory.
func NewCache(capacity int) *Cache {
	return &Cache{
		entries: make(map[string]*entry, capacity),
		evict:   NewLRUEviction(capacity),
	}
}

// Put stores an index for the given file path.
func (c *Cache) Put(path string, idx *Index) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if evicted := c.evict.Touch(path); evicted != "" {
		delete(c.entries, evicted)
	}
	c.entries[path] = &entry{index: idx, modTime: info.ModTime()}
	return nil
}

// Get retrieves a cached index. Returns nil, false if not present or stale.
func (c *Cache) Get(path string) (*Index, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.entries[path]
	if !ok {
		return nil, false
	}
	info, err := os.Stat(path)
	if err != nil || info.ModTime().After(e.modTime) {
		c.evict.Remove(path)
		delete(c.entries, path)
		return nil, false
	}
	c.evict.Touch(path)
	return e.index, true
}

// Invalidate removes a cached entry for the given path.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict.Remove(path)
	delete(c.entries, path)
}

// IsStale reports whether the cached entry for path is outdated.
func (c *Cache) IsStale(path string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[path]
	if !ok {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	return info.ModTime().After(e.modTime)
}

// Len returns the number of entries currently in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
