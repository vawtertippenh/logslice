package index

import (
	"testing"
	"time"
)

func makeStats(entries, hits, misses int, first, last time.Time) Stats {
	return Stats{
		Entries:      entries,
		SampleRate:   1,
		FirstTime:    first,
		LastTime:     last,
		CacheHits:    hits,
		CacheMisses:  misses,
		BytesCovered: 1024,
	}
}

func TestStatsDuration(t *testing.T) {
	now := time.Now()
	later := now.Add(5 * time.Minute)
	s := makeStats(10, 0, 0, now, later)
	if got := s.Duration(); got != 5*time.Minute {
		t.Errorf("Duration() = %v, want %v", got, 5*time.Minute)
	}
}

func TestStatsDurationZero(t *testing.T) {
	s := Stats{}
	if got := s.Duration(); got != 0 {
		t.Errorf("Duration() = %v, want 0", got)
	}
}

func TestStatsIsEmpty(t *testing.T) {
	s := Stats{}
	if !s.IsEmpty() {
		t.Error("IsEmpty() = false, want true for zero stats")
	}
	s.Entries = 5
	if s.IsEmpty() {
		t.Error("IsEmpty() = true, want false when entries > 0")
	}
}

func TestStatsCacheHitRate(t *testing.T) {
	cases := []struct {
		hits, misses int
		want         float64
	}{
		{0, 0, 0},
		{10, 0, 1.0},
		{0, 10, 0.0},
		{3, 7, 0.3},
	}
	for _, tc := range cases {
		s := Stats{CacheHits: tc.hits, CacheMisses: tc.misses}
		got := s.CacheHitRate()
		if got != tc.want {
			t.Errorf("CacheHitRate() hits=%d misses=%d = %v, want %v",
				tc.hits, tc.misses, got, tc.want)
		}
	}
}
