package index

import "time"

// Stats holds runtime statistics about an index.
type Stats struct {
	// Entries is the total number of index entries.
	Entries int

	// SampleRate is the sampling interval used when building the index.
	SampleRate int

	// FirstTime is the timestamp of the first indexed entry.
	FirstTime time.Time

	// LastTime is the timestamp of the last indexed entry.
	LastTime time.Time

	// CacheHits is the number of times a cached index was returned.
	CacheHits int

	// CacheMisses is the number of times the cache did not have an entry.
	CacheMisses int

	// BytesCovered is the total number of bytes spanned by the index.
	BytesCovered int64
}

// Duration returns the time span covered by the index.
func (s Stats) Duration() time.Duration {
	if s.FirstTime.IsZero() || s.LastTime.IsZero() {
		return 0
	}
	return s.LastTime.Sub(s.FirstTime)
}

// IsEmpty reports whether the index contains no entries.
func (s Stats) IsEmpty() bool {
	return s.Entries == 0
}

// CacheHitRate returns the ratio of cache hits to total cache lookups.
// Returns 0 if no lookups have been made.
func (s Stats) CacheHitRate() float64 {
	total := s.CacheHits + s.CacheMisses
	if total == 0 {
		return 0
	}
	return float64(s.CacheHits) / float64(total)
}
