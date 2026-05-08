package multiline

import (
	"errors"
	"regexp"
	"strings"
)

// Options controls aggregator behaviour.
type Options struct {
	// StartPattern is a regular expression that matches the first line of a
	// new log record.  Lines that do NOT match are treated as continuations
	// of the previous record.
	StartPattern string

	// MaxLines caps how many physical lines are joined into one record.
	// Zero means no limit.
	MaxLines int

	// Separator is placed between joined lines (default: "\n").
	Separator string
}

// Aggregator joins continuation lines into complete log records.
type Aggregator struct {
	start    *regexp.Regexp
	maxLines int
	sep      string
	buf      []string
}

// New creates an Aggregator from opts.  StartPattern must be a valid regexp.
func New(opts Options) (*Aggregator, error) {
	if opts.StartPattern == "" {
		return nil, errors.New("multiline: StartPattern must not be empty")
	}
	re, err := regexp.Compile(opts.StartPattern)
	if err != nil {
		return nil, err
	}
	sep := opts.Separator
	if sep == "" {
		sep = "\n"
	}
	return &Aggregator{
		start:    re,
		maxLines: opts.MaxLines,
		sep:      sep,
	}, nil
}

// Feed accepts a raw line.  If feeding this line completes a previous record,
// that record is returned with ok == true.  Otherwise ok is false.
func (a *Aggregator) Feed(line string) (record string, ok bool) {
	isStart := a.start.MatchString(line)

	if isStart && len(a.buf) > 0 {
		record = strings.Join(a.buf, a.sep)
		a.buf = []string{line}
		return record, true
	}

	a.buf = append(a.buf, line)

	// Force flush when the buffer hits the line cap.
	if a.maxLines > 0 && len(a.buf) >= a.maxLines {
		return a.Flush()
	}

	return "", false
}

// Flush returns any buffered lines as a completed record and resets the buffer.
// Call this after the input is exhausted to retrieve the final record.
func (a *Aggregator) Flush() (record string, ok bool) {
	if len(a.buf) == 0 {
		return "", false
	}
	record = strings.Join(a.buf, a.sep)
	a.buf = a.buf[:0]
	return record, true
}

// Reset discards any buffered state.
func (a *Aggregator) Reset() {
	a.buf = a.buf[:0]
}
