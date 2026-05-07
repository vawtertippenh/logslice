// Package sampler provides line-based and byte-based sampling for large log files.
// It allows callers to read every Nth line or limit output to a maximum number of lines.
package sampler

import (
	"errors"
)

// Options configures the Sampler behaviour.
type Options struct {
	// Every emits one line for every N lines seen. Must be >= 1.
	Every int
	// MaxLines caps the total number of lines emitted. Zero means unlimited.
	MaxLines int
}

// Sampler decides whether a given line should be emitted based on sampling rules.
type Sampler struct {
	opts    Options
	seen    int
	emitted int
}

// New creates a Sampler with the given Options.
// Returns an error if Every is less than 1.
func New(opts Options) (*Sampler, error) {
	if opts.Every < 1 {
		return nil, errors.New("sampler: Every must be >= 1")
	}
	return &Sampler{opts: opts}, nil
}

// ShouldEmit returns true if the current line should be included in output.
// It must be called exactly once per line in order.
func (s *Sampler) ShouldEmit() bool {
	s.seen++
	if s.opts.MaxLines > 0 && s.emitted >= s.opts.MaxLines {
		return false
	}
	if s.seen%s.opts.Every != 0 {
		return false
	}
	s.emitted++
	return true
}

// Reset clears all counters, allowing the Sampler to be reused.
func (s *Sampler) Reset() {
	s.seen = 0
	s.emitted = 0
}

// Seen returns the total number of lines evaluated.
func (s *Sampler) Seen() int { return s.seen }

// Emitted returns the total number of lines emitted.
func (s *Sampler) Emitted() int { return s.emitted }
