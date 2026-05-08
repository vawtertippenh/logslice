// Package multiline provides support for grouping multi-line log entries
// into a single logical record before filtering or output.
//
// Some log formats (e.g. Java stack traces, Python tracebacks) emit a single
// logical event across multiple physical lines. A Joiner accumulates lines
// until a configurable start-of-record pattern is matched, then emits the
// previously accumulated group as one joined string.
package multiline

import (
	"regexp"
	"strings"
)

// Options controls Joiner behaviour.
type Options struct {
	// StartPattern is a regular expression that marks the first line of a new
	// logical record. Every line that matches begins a fresh group.
	StartPattern *regexp.Regexp

	// Separator is placed between physical lines when joining. Defaults to "\n".
	Separator string
}

// Joiner accumulates physical log lines and emits logical records.
type Joiner struct {
	opts    Options
	buf     []string
	pending string
}

// New creates a Joiner with the provided options.
// StartPattern must be non-nil.
func New(opts Options) *Joiner {
	if opts.Separator == "" {
		opts.Separator = "\n"
	}
	return &Joiner{opts: opts}
}

// Add feeds the next physical line to the Joiner.
// If the line starts a new record, the previously buffered record is returned
// together with ok=true. Otherwise ok is false and the caller should continue
// feeding lines.
func (j *Joiner) Add(line string) (record string, ok bool) {
	if j.opts.StartPattern.MatchString(line) {
		if len(j.buf) > 0 {
			record = strings.Join(j.buf, j.opts.Separator)
			ok = true
		}
		j.buf = []string{line}
		return
	}
	j.buf = append(j.buf, line)
	return
}

// Flush returns any remaining buffered lines as a final record.
// It should be called after the input is exhausted.
func (j *Joiner) Flush() (record string, ok bool) {
	if len(j.buf) == 0 {
		return
	}
	record = strings.Join(j.buf, j.opts.Separator)
	ok = true
	j.buf = nil
	return
}

// Reset clears all buffered state, ready for a new input stream.
func (j *Joiner) Reset() {
	j.buf = nil
}
