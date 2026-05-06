// Package ratelimit provides a token-bucket rate limiter for controlling
// the throughput of log line output, useful when streaming large files.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter that controls the number of
// lines (or bytes) emitted per second.
type Limiter struct {
	mu       sync.Mutex
	rate     float64 // tokens per second
	burst    float64 // maximum burst size
	tokens   float64
	lastTick time.Time
	clock    func() time.Time
}

// Options configures the rate limiter.
type Options struct {
	// Rate is the number of tokens replenished per second.
	Rate float64
	// Burst is the maximum number of tokens that can accumulate.
	Burst float64
}

// New returns a new Limiter with the given options.
// If Burst is zero it defaults to Rate.
func New(opts Options) *Limiter {
	if opts.Burst <= 0 {
		opts.Burst = opts.Rate
	}
	return &Limiter{
		rate:     opts.Rate,
		burst:    opts.Burst,
		tokens:   opts.Burst,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Wait blocks until one token is available, then consumes it.
func (l *Limiter) Wait() {
	l.WaitN(1)
}

// WaitN blocks until n tokens are available, then consumes them.
func (l *Limiter) WaitN(n float64) {
	for {
		l.mu.Lock()
		l.refill()
		if l.tokens >= n {
			l.tokens -= n
			l.mu.Unlock()
			return
		}
		// Calculate how long to wait for enough tokens.
		need := n - l.tokens
		wait := time.Duration(need/l.rate*float64(time.Second))
		l.mu.Unlock()
		time.Sleep(wait)
	}
}

// Available returns the current number of available tokens without blocking.
func (l *Limiter) Available() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	return l.tokens
}

// refill adds tokens based on elapsed time. Must be called with l.mu held.
func (l *Limiter) refill() {
	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}
	l.lastTick = now
}
