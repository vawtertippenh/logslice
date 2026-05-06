package index

import (
	"testing"
	"time"
)

func TestTTLCachePutAndGet(t *testing.T) {
	c := NewTTLCache(time.Minute)
	c.Put("key1", "value1")

	v, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to be present")
	}
	if v.(string) != "value1" {
		t.Fatalf("expected value1, got %v", v)
	}
}

func TestTTLCacheGetMissing(t *testing.T) {
	c := NewTTLCache(time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestTTLCacheExpiry(t *testing.T) {
	c := NewTTLCache(10 * time.Millisecond)
	c.Put("expires", 42)

	// Should be present immediately.
	_, ok := c.Get("expires")
	if !ok {
		t.Fatal("expected entry to be present before expiry")
	}

	time.Sleep(20 * time.Millisecond)

	_, ok = c.Get("expires")
	if ok {
		t.Fatal("expected entry to have expired")
	}
}

func TestTTLCacheDelete(t *testing.T) {
	c := NewTTLCache(time.Minute)
	c.Put("del", "v")
	c.Delete("del")

	_, ok := c.Get("del")
	if ok {
		t.Fatal("expected deleted key to be absent")
	}
}

func TestTTLCachePurge(t *testing.T) {
	c := NewTTLCache(10 * time.Millisecond)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	time.Sleep(20 * time.Millisecond)

	// Add a fresh entry that should not be purged.
	c.Put("fresh", 4)

	removed := c.Purge()
	if removed != 3 {
		t.Fatalf("expected 3 removed, got %d", removed)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", c.Len())
	}
}

func TestTTLCacheLen(t *testing.T) {
	c := NewTTLCache(time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Put("x", 1)
	c.Put("y", 2)
	if c.Len() != 2 {
		t.Fatalf("expected len 2, got %d", c.Len())
	}
}

func TestTTLEntryIsExpired(t *testing.T) {
	past := TTLEntry{value: "v", expiresAt: time.Now().Add(-time.Second)}
	if !past.IsExpired(time.Now()) {
		t.Fatal("expected past entry to be expired")
	}

	future := TTLEntry{value: "v", expiresAt: time.Now().Add(time.Hour)}
	if future.IsExpired(time.Now()) {
		t.Fatal("expected future entry to not be expired")
	}
}
