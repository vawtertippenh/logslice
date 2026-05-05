package index

import (
	"testing"
	"time"
)

func buildIndex(pairs [][2]int64) *Index {
	idx := New()
	for _, p := range pairs {
		idx.addRaw(p[0], p[1])
	}
	return idx
}

func TestMergerBasic(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
	sec := int64(time.Second)

	a := buildIndex([][2]int64{{t0, 0}, {t0 + 2*sec, 200}})
	b := buildIndex([][2]int64{{t0 + sec, 100}, {t0 + 3*sec, 300}})

	m := NewMerger(a, b)
	out := m.Merge()

	if out.Len() != 4 {
		t.Fatalf("expected 4 entries, got %d", out.Len())
	}

	for i := 0; i < out.Len()-1; i++ {
		if !out.EntryAt(i).Time.Before(out.EntryAt(i + 1).Time) {
			t.Errorf("entries not sorted at position %d", i)
		}
	}
}

func TestMergerDeduplicatesOffsets(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
	sec := int64(time.Second)

	a := buildIndex([][2]int64{{t0, 0}, {t0 + sec, 100}})
	b := buildIndex([][2]int64{{t0, 0}, {t0 + 2*sec, 200}})

	out := NewMerger(a, b).Merge()

	if out.Len() != 3 {
		t.Fatalf("expected 3 entries after dedup, got %d", out.Len())
	}
}

func TestMergerNilSource(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
	a := buildIndex([][2]int64{{t0, 0}})

	out := NewMerger(a, nil).Merge()
	if out.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", out.Len())
	}
}

func TestMergerEmpty(t *testing.T) {
	out := NewMerger().Merge()
	if out.Len() != 0 {
		t.Fatalf("expected empty index, got %d entries", out.Len())
	}
}
