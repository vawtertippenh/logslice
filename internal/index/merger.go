package index

import "sort"

// Merger combines multiple sorted Index slices into a single deduplicated,
// sorted Index. This is useful when building a unified index from several
// independently indexed log segments.
type Merger struct {
	sources []*Index
}

// NewMerger returns a Merger that will combine the provided indexes.
func NewMerger(sources ...*Index) *Merger {
	return &Merger{sources: sources}
}

// Merge combines all source indexes into a new Index.
// Entries are sorted by timestamp and deduplicated by byte offset.
func (m *Merger) Merge() *Index {
	type entry struct {
		t  int64 // unix nano
		off int64
	}

	seen := make(map[int64]struct{})
	var entries []entry

	for _, idx := range m.sources {
		if idx == nil {
			continue
		}
		for i := 0; i < idx.Len(); i++ {
			e := idx.EntryAt(i)
			if _, dup := seen[e.Offset]; dup {
				continue
			}
			seen[e.Offset] = struct{}{}
			entries = append(entries, entry{t: e.Time.UnixNano(), off: e.Offset})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].t != entries[j].t {
			return entries[i].t < entries[j].t
		}
		return entries[i].off < entries[j].off
	})

	out := New()
	for _, e := range entries {
		out.addRaw(e.t, e.off)
	}
	return out
}
