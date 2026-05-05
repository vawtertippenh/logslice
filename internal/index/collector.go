package index

import "sync/atomic"

// Collector accumulates cache hit/miss counts in a thread-safe manner
// and can produce a Stats snapshot when combined with an Index.
type Collector struct {
	hits   int64
	misses int64
}

// RecordHit increments the cache hit counter.
func (c *Collector) RecordHit() {
	atomic.AddInt64(&c.hits, 1)
}

// RecordMiss increments the cache miss counter.
func (c *Collector) RecordMiss() {
	atomic.AddInt64(&c.misses, 1)
}

// Hits returns the current hit count.
func (c *Collector) Hits() int {
	return int(atomic.LoadInt64(&c.hits))
}

// Misses returns the current miss count.
func (c *Collector) Misses() int {
	return int(atomic.LoadInt64(&c.misses))
}

// Reset zeroes all counters.
func (c *Collector) Reset() {
	atomic.StoreInt64(&c.hits, 0)
	atomic.StoreInt64(&c.misses, 0)
}

// StatsFrom builds a Stats value from this Collector and the provided Index.
func (c *Collector) StatsFrom(idx *Index) Stats {
	n := idx.Len()
	var first, last entry
	if n > 0 {
		first = idx.entries[0]
		last = idx.entries[n-1]
	}
	var bytesCovered int64
	if n > 0 {
		bytesCovered = last.Offset - first.Offset
	}
	return Stats{
		Entries:      n,
		SampleRate:   idx.sampleRate,
		FirstTime:    first.Time,
		LastTime:     last.Time,
		CacheHits:    c.Hits(),
		CacheMisses:  c.Misses(),
		BytesCovered: bytesCovered,
	}
}
