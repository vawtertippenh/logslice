// Package truncate provides utilities for truncating long log lines
// to a configurable maximum byte length, preserving valid UTF-8 boundaries.
package truncate

import (
	"unicode/utf8"
)

// DefaultMaxBytes is the default maximum line length in bytes.
const DefaultMaxBytes = 4096

// Suffix is appended to truncated lines to indicate truncation.
const Suffix = "..."

// Options configures the Truncator.
type Options struct {
	// MaxBytes is the maximum number of bytes allowed per line.
	// Lines exceeding this length will be truncated. Must be > len(Suffix).
	MaxBytes int

	// AppendSuffix controls whether the truncation suffix is appended.
	AppendSuffix bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxBytes:     DefaultMaxBytes,
		AppendSuffix: true,
	}
}

// Truncator truncates byte slices to a maximum length while
// respecting UTF-8 character boundaries.
type Truncator struct {
	opts Options
}

// New creates a new Truncator with the given options.
// If opts.MaxBytes is <= 0, DefaultMaxBytes is used.
func New(opts Options) *Truncator {
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = DefaultMaxBytes
	}
	return &Truncator{opts: opts}
}

// Truncate returns the input slice truncated to the configured maximum
// byte length. If the input fits within the limit, it is returned unchanged.
// The returned slice may share memory with the input.
func (t *Truncator) Truncate(line []byte) []byte {
	if len(line) <= t.opts.MaxBytes {
		return line
	}

	limit := t.opts.MaxBytes
	if t.opts.AppendSuffix {
		limit -= len(Suffix)
		if limit < 0 {
			limit = 0
		}
	}

	// Walk back to a valid UTF-8 boundary.
	cut := limit
	for cut > 0 && !utf8.RuneStart(line[cut]) {
		cut--
	}

	if t.opts.AppendSuffix {
		out := make([]byte, cut, cut+len(Suffix))
		copy(out, line[:cut])
		out = append(out, Suffix...)
		return out
	}

	out := make([]byte, cut)
	copy(out, line[:cut])
	return out
}

// TruncateString is a convenience wrapper around Truncate for strings.
func (t *Truncator) TruncateString(line string) string {
	return string(t.Truncate([]byte(line)))
}
