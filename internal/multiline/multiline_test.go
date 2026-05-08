package multiline

import (
	"strings"
	"testing"
)

func TestNewBadPattern(t *testing.T) {
	_, err := New(Options{StartPattern: "["})
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestNewEmptyPattern(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestFeedSingleRecord(t *testing.T) {
	agg := mustNew(t, `^\d{4}-`)
	lines := []string{
		"2024-01-01 start",
		"  continuation",
		"  more",
		"2024-01-02 next",
	}

	var records []string
	for _, l := range lines {
		if rec, ok := agg.Feed(l); ok {
			records = append(records, rec)
		}
	}
	if rec, ok := agg.Flush(); ok {
		records = append(records, rec)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if !strings.Contains(records[0], "continuation") {
		t.Errorf("first record missing continuation: %q", records[0])
	}
}

func TestFlushEmpty(t *testing.T) {
	agg := mustNew(t, `^START`)
	_, ok := agg.Flush()
	if ok {
		t.Fatal("expected ok=false on empty flush")
	}
}

func TestMaxLines(t *testing.T) {
	agg, _ := New(Options{StartPattern: `^START`, MaxLines: 3})
	agg.Feed("START line")
	agg.Feed("cont 1")
	rec, ok := agg.Feed("cont 2") // hits cap
	if !ok {
		t.Fatal("expected flush at MaxLines")
	}
	if !strings.Contains(rec, "cont 2") {
		t.Errorf("record missing last line: %q", rec)
	}
}

func TestReset(t *testing.T) {
	agg := mustNew(t, `^START`)
	agg.Feed("START foo")
	agg.Reset()
	_, ok := agg.Flush()
	if ok {
		t.Fatal("expected empty buffer after Reset")
	}
}

func TestCustomSeparator(t *testing.T) {
	agg, _ := New(Options{StartPattern: `^LOG`, Separator: " | "})
	agg.Feed("LOG begin")
	agg.Feed("detail")
	rec, _ := agg.Flush()
	if !strings.Contains(rec, " | ") {
		t.Errorf("expected custom separator in %q", rec)
	}
}

func mustNew(t *testing.T, pattern string) *Aggregator {
	t.Helper()
	agg, err := New(Options{StartPattern: pattern})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return agg
}
