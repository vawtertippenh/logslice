// Package index provides byte-offset indexing for large log files,
// enabling fast seeking to time-range boundaries without full scans.
package index

import (
	"time"
)

// Entry represents a single index entry mapping a timestamp to a byte offset.
type Entry struct {
	Timestamp time.Time
	Offset    int64
}

// Index holds a sorted list of log file entries for binary search.
type Index struct {
	entries []Entry
}

// New creates an empty Index.
func New() *Index {
	return &Index{}
}

// Add appends a new entry to the index. Entries should be added in
// chronological order for binary search to work correctly.
func (idx *Index) Add(ts time.Time, offset int64) {
	idx.entries = append(idx.entries, Entry{Timestamp: ts, Offset: offset})
}

// Len returns the number of entries in the index.
func (idx *Index) Len() int {
	return len(idx.entries)
}

// FindStart returns the byte offset of the first entry whose timestamp
// is >= start. Returns 0 if no such entry exists or index is empty.
func (idx *Index) FindStart(start time.Time) int64 {
	for _, e := range idx.entries {
		if !e.Timestamp.Before(start) {
			return e.Offset
		}
	}
	return 0
}

// FindEnd returns the byte offset just past the last entry whose timestamp
// is <= end. Returns -1 to indicate read until EOF.
func (idx *Index) FindEnd(end time.Time) int64 {
	last := int64(-1)
	for _, e := range idx.entries {
		if !e.Timestamp.After(end) {
			last = e.Offset
		} else {
			break
		}
	}
	return last
}

// Entries returns a copy of all index entries.
func (idx *Index) Entries() []Entry {
	result := make([]Entry, len(idx.entries))
	copy(result, idx.entries)
	return result
}
