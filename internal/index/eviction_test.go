package index

import (
	"fmt"
	"testing"
)

func TestLRUEvictionTouch(t *testing.T) {
	e := NewLRUEviction(3)

	if evicted := e.Touch("a"); evicted != "" {
		t.Errorf("expected no eviction, got %q", evicted)
	}
	e.Touch("b")
	e.Touch("c")

	if e.Len() != 3 {
		t.Errorf("expected len 3, got %d", e.Len())
	}

	// Adding a 4th entry should evict the oldest ("a").
	evicted := e.Touch("d")
	if evicted != "a" {
		t.Errorf("expected eviction of \"a\", got %q", evicted)
	}
	if e.Len() != 3 {
		t.Errorf("expected len 3 after eviction, got %d", e.Len())
	}
}

func TestLRUEvictionTouchExisting(t *testing.T) {
	e := NewLRUEviction(3)
	e.Touch("a")
	e.Touch("b")
	e.Touch("c")

	// Re-touch "a" to make it most-recently-used.
	e.Touch("a")

	// Next eviction should remove "b", not "a".
	evicted := e.Touch("d")
	if evicted != "b" {
		t.Errorf("expected eviction of \"b\", got %q", evicted)
	}
}

func TestLRUEvictionRemove(t *testing.T) {
	e := NewLRUEviction(3)
	e.Touch("a")
	e.Touch("b")
	e.Remove("a")

	if e.Len() != 1 {
		t.Errorf("expected len 1, got %d", e.Len())
	}

	// Removing a non-existent key should not panic.
	e.Remove("z")
}

func TestLRUEvictionDefaultCapacity(t *testing.T) {
	e := NewLRUEviction(0)
	for i := 0; i < 8; i++ {
		e.Touch(fmt.Sprintf("key%d", i))
	}
	if e.Len() != 8 {
		t.Errorf("expected len 8, got %d", e.Len())
	}
	// 9th entry triggers eviction.
	evicted := e.Touch("overflow")
	if evicted == "" {
		t.Error("expected an eviction with default capacity of 8")
	}
}

func TestLRUEvictionCapacityOne(t *testing.T) {
	e := NewLRUEviction(1)
	e.Touch("first")
	evicted := e.Touch("second")
	if evicted != "first" {
		t.Errorf("expected \"first\" evicted, got %q", evicted)
	}
	if e.Len() != 1 {
		t.Errorf("expected len 1, got %d", e.Len())
	}
}
