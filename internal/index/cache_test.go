package index

import (
	"testing"
	"time"
)

func TestCachePutAndGet(t *testing.T) {
	c := NewCache()
	idx := New()

	c.Put("/var/log/app.log", idx, 1024)

	entry, ok := c.Get("/var/log/app.log")
	if !ok {
		t.Fatal("expected entry to exist in cache")
	}
	if entry.Index != idx {
		t.Error("cached index does not match stored index")
	}
	if entry.FileSize != 1024 {
		t.Errorf("expected file size 1024, got %d", entry.FileSize)
	}
	if entry.FilePath != "/var/log/app.log" {
		t.Errorf("unexpected file path: %s", entry.FilePath)
	}
	if entry.BuiltAt.IsZero() {
		t.Error("BuiltAt should not be zero")
	}
	if entry.BuiltAt.After(time.Now()) {
		t.Error("BuiltAt should not be in the future")
	}
}

func TestCacheGetMissing(t *testing.T) {
	c := NewCache()
	_, ok := c.Get("/nonexistent.log")
	if ok {
		t.Error("expected cache miss for unknown path")
	}
}

func TestCacheInvalidate(t *testing.T) {
	c := NewCache()
	c.Put("/var/log/app.log", New(), 512)
	c.Invalidate("/var/log/app.log")

	_, ok := c.Get("/var/log/app.log")
	if ok {
		t.Error("expected entry to be removed after invalidation")
	}
}

func TestCacheIsStale(t *testing.T) {
	c := NewCache()
	c.Put("/var/log/app.log", New(), 1000)

	if c.IsStale("/var/log/app.log", 1000) {
		t.Error("entry with same size should not be stale")
	}
	if !c.IsStale("/var/log/app.log", 2000) {
		t.Error("entry with different size should be stale")
	}
	if !c.IsStale("/missing.log", 100) {
		t.Error("missing entry should be considered stale")
	}
}

func TestCacheLen(t *testing.T) {
	c := NewCache()
	if c.Len() != 0 {
		t.Errorf("expected length 0, got %d", c.Len())
	}

	c.Put("/a.log", New(), 100)
	c.Put("/b.log", New(), 200)
	if c.Len() != 2 {
		t.Errorf("expected length 2, got %d", c.Len())
	}

	c.Invalidate("/a.log")
	if c.Len() != 1 {
		t.Errorf("expected length 1 after invalidation, got %d", c.Len())
	}
}
