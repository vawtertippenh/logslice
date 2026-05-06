package ratelimit

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiterAvailableInitial(t *testing.T) {
	l := New(Options{Rate: 10, Burst: 10})
	got := l.Available()
	if got != 10 {
		t.Fatalf("expected 10 initial tokens, got %v", got)
	}
}

func TestLimiterBurstDefaultsToRate(t *testing.T) {
	l := New(Options{Rate: 5})
	if l.burst != 5 {
		t.Fatalf("expected burst=5, got %v", l.burst)
	}
}

func TestLimiterWaitConsumesToken(t *testing.T) {
	l := New(Options{Rate: 100, Burst: 10})
	l.Wait()
	got := l.Available()
	// After one Wait the token count should have dropped by (at least) 1
	// but refill may have added a small fraction back.
	if got > 9.1 {
		t.Fatalf("expected tokens < 9.1 after Wait, got %v", got)
	}
}

func TestLimiterRefillOverTime(t *testing.T) {
	now := time.Now()
	l := New(Options{Rate: 10, Burst: 10})
	// Drain all tokens.
	l.mu.Lock()
	l.tokens = 0
	l.lastTick = now
	l.mu.Unlock()

	// Advance the clock by 1 second via the internal clock override.
	l.mu.Lock()
	l.clock = func() time.Time { return now.Add(time.Second) }
	l.mu.Unlock()

	got := l.Available()
	if got < 9.9 || got > 10.0 {
		t.Fatalf("expected ~10 tokens after 1s refill, got %v", got)
	}
}

func TestLimiterBurstCap(t *testing.T) {
	now := time.Now()
	l := New(Options{Rate: 10, Burst: 5})
	l.mu.Lock()
	l.tokens = 0
	l.lastTick = now
	l.clock = func() time.Time { return now.Add(10 * time.Second) }
	l.mu.Unlock()

	got := l.Available()
	if got != 5 {
		t.Fatalf("expected tokens capped at burst=5, got %v", got)
	}
}

func TestLimiterConcurrentWait(t *testing.T) {
	l := New(Options{Rate: 1000, Burst: 1000})
	var count int64
	workers := 20
	done := make(chan struct{})
	for i := 0; i < workers; i++ {
		go func() {
			l.Wait()
			atomic.AddInt64(&count, 1)
			done <- struct{}{}
		}()
	}
	timeout := time.After(2 * time.Second)
	for i := 0; i < workers; i++ {
		select {
		case <-done:
		case <-timeout:
			t.Fatalf("timed out waiting for workers; completed %d/%d", atomic.LoadInt64(&count), workers)
		}
	}
}
