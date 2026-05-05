package index

import (
	"sync"
	"testing"
	"time"
)

func TestCollectorRecordAndRead(t *testing.T) {
	var c Collector
	c.RecordHit()
	c.RecordHit()
	c.RecordMiss()
	if got := c.Hits(); got != 2 {
		t.Errorf("Hits() = %d, want 2", got)
	}
	if got := c.Misses(); got != 1 {
		t.Errorf("Misses() = %d, want 1", got)
	}
}

func TestCollectorReset(t *testing.T) {
	var c Collector
	c.RecordHit()
	c.RecordMiss()
	c.Reset()
	if c.Hits() != 0 || c.Misses() != 0 {
		t.Error("Reset() did not zero counters")
	}
}

func TestCollectorConcurrent(t *testing.T) {
	var c Collector
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.RecordHit()
			c.RecordMiss()
		}()
	}
	wg.Wait()
	if c.Hits() != 100 {
		t.Errorf("Hits() = %d, want 100", c.Hits())
	}
	if c.Misses() != 100 {
		t.Errorf("Misses() = %d, want 100", c.Misses())
	}
}

func TestCollectorStatsFrom(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := New(2)
	idx.Add(ts, 0)
	idx.Add(ts.Add(time.Minute), 512)
	idx.Add(ts.Add(2*time.Minute), 1024)

	var c Collector
	c.RecordHit()
	c.RecordMiss()

	s := c.StatsFrom(idx)
	if s.Entries != 2 {
		t.Errorf("Entries = %d, want 2", s.Entries)
	}
	if s.CacheHits != 1 {
		t.Errorf("CacheHits = %d, want 1", s.CacheHits)
	}
	if s.CacheMisses != 1 {
		t.Errorf("CacheMisses = %d, want 1", s.CacheMisses)
	}
	if s.BytesCovered <= 0 {
		t.Errorf("BytesCovered = %d, want > 0", s.BytesCovered)
	}
}
