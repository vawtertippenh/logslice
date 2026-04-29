package index

import (
	"testing"
	"time"
)

func makeTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestIndexAddAndLen(t *testing.T) {
	idx := New()
	if idx.Len() != 0 {
		t.Fatalf("expected 0, got %d", idx.Len())
	}
	idx.Add(makeTime("2024-01-01T10:00:00Z"), 0)
	idx.Add(makeTime("2024-01-01T10:05:00Z"), 512)
	if idx.Len() != 2 {
		t.Fatalf("expected 2, got %d", idx.Len())
	}
}

func TestIndexFindStart(t *testing.T) {
	idx := New()
	idx.Add(makeTime("2024-01-01T10:00:00Z"), 0)
	idx.Add(makeTime("2024-01-01T10:05:00Z"), 512)
	idx.Add(makeTime("2024-01-01T10:10:00Z"), 1024)

	offset := idx.FindStart(makeTime("2024-01-01T10:05:00Z"))
	if offset != 512 {
		t.Errorf("expected 512, got %d", offset)
	}

	offset = idx.FindStart(makeTime("2024-01-01T09:00:00Z"))
	if offset != 0 {
		t.Errorf("expected 0 for early start, got %d", offset)
	}
}

func TestIndexFindEnd(t *testing.T) {
	idx := New()
	idx.Add(makeTime("2024-01-01T10:00:00Z"), 0)
	idx.Add(makeTime("2024-01-01T10:05:00Z"), 512)
	idx.Add(makeTime("2024-01-01T10:10:00Z"), 1024)

	offset := idx.FindEnd(makeTime("2024-01-01T10:05:00Z"))
	if offset != 512 {
		t.Errorf("expected 512, got %d", offset)
	}

	offset = idx.FindEnd(makeTime("2024-01-01T12:00:00Z"))
	if offset != 1024 {
		t.Errorf("expected 1024, got %d", offset)
	}
}

func TestIndexFindEndNoMatch(t *testing.T) {
	idx := New()
	idx.Add(makeTime("2024-01-01T10:00:00Z"), 0)

	offset := idx.FindEnd(makeTime("2024-01-01T09:00:00Z"))
	if offset != -1 {
		t.Errorf("expected -1, got %d", offset)
	}
}

func TestIndexEntries(t *testing.T) {
	idx := New()
	idx.Add(makeTime("2024-01-01T10:00:00Z"), 0)
	idx.Add(makeTime("2024-01-01T10:05:00Z"), 512)

	entries := idx.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Offset != 0 || entries[1].Offset != 512 {
		t.Errorf("unexpected entry offsets: %+v", entries)
	}
}
