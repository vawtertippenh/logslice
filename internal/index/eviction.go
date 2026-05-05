package index

import (
	"container/list"
	"sync"
)

// LRUEviction implements a least-recently-used eviction policy for the index cache.
// It tracks access order and evicts the oldest entry when capacity is exceeded.
type LRUEviction struct {
	mu       sync.Mutex
	capacity int
	order    *list.List
	keys     map[string]*list.Element
}

// NewLRUEviction creates a new LRUEviction with the given capacity.
func NewLRUEviction(capacity int) *LRUEviction {
	if capacity <= 0 {
		capacity = 8
	}
	return &LRUEviction{
		capacity: capacity,
		order:    list.New(),
		keys:     make(map[string]*list.Element, capacity),
	}
}

// Touch records an access to the given key, moving it to the front.
// Returns the evicted key if capacity was exceeded, or an empty string.
func (e *LRUEviction) Touch(key string) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if elem, ok := e.keys[key]; ok {
		e.order.MoveToFront(elem)
		return ""
	}

	elem := e.order.PushFront(key)
	e.keys[key] = elem

	if e.order.Len() > e.capacity {
		return e.evictOldest()
	}
	return ""
}

// Remove explicitly removes a key from the eviction tracker.
func (e *LRUEviction) Remove(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if elem, ok := e.keys[key]; ok {
		e.order.Remove(elem)
		delete(e.keys, key)
	}
}

// Len returns the number of tracked keys.
func (e *LRUEviction) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.order.Len()
}

// evictOldest removes and returns the least-recently-used key.
// Caller must hold e.mu.
func (e *LRUEviction) evictOldest() string {
	back := e.order.Back()
	if back == nil {
		return ""
	}
	key := back.Value.(string)
	e.order.Remove(back)
	delete(e.keys, key)
	return key
}
