package index

import (
	"sync"
	"time"
)

// CacheEntry holds a cached index along with metadata about when it was built.
type CacheEntry struct {
	Index     *Index
	BuiltAt   time.Time
	FileSize  int64
	FilePath  string
}

// Cache stores built indexes keyed by file path, allowing reuse across
// multiple queries on the same file without rebuilding the index.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
}

// NewCache creates an empty index cache.
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry),
	}
}

// Get retrieves a cached index entry for the given file path.
// Returns nil, false if no entry exists.
func (c *Cache) Get(filePath string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[filePath]
	return entry, ok
}

// Put stores an index entry in the cache for the given file path.
func (c *Cache) Put(filePath string, idx *Index, fileSize int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[filePath] = &CacheEntry{
		Index:    idx,
		BuiltAt:  time.Now(),
		FileSize: fileSize,
		FilePath: filePath,
	}
}

// Invalidate removes the cached entry for the given file path.
func (c *Cache) Invalidate(filePath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, filePath)
}

// IsStale returns true if the cached entry's recorded file size differs
// from the provided current size, indicating the file has changed.
func (c *Cache) IsStale(filePath string, currentSize int64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[filePath]
	if !ok {
		return true
	}
	return entry.FileSize != currentSize
}

// Len returns the number of entries currently in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
