package multiline

import (
	"bufio"
	"io"
)

// Scanner wraps an Aggregator and an io.Reader to present a line-oriented
// interface that yields complete (possibly multi-line) log records one at a
// time, similar to bufio.Scanner.
type Scanner struct {
	src  *bufio.Scanner
	agg  *Aggregator
	rec  string
	err  error
	done bool
}

// NewScanner returns a Scanner that reads from r and aggregates lines
// according to opts.
func NewScanner(r io.Reader, opts Options) (*Scanner, error) {
	agg, err := New(opts)
	if err != nil {
		return nil, err
	}
	return &Scanner{
		src: bufio.NewScanner(r),
		agg: agg,
	}, nil
}

// Scan advances to the next complete record.  It returns true if a record is
// available via Text, false when the input is exhausted or an error occurs.
func (s *Scanner) Scan() bool {
	if s.done {
		return false
	}
	for s.src.Scan() {
		line := s.src.Text()
		if rec, ok := s.agg.Feed(line); ok {
			s.rec = rec
			return true
		}
	}
	if err := s.src.Err(); err != nil {
		s.err = err
		s.done = true
		return false
	}
	// Flush the last buffered record.
	s.done = true
	if rec, ok := s.agg.Flush(); ok {
		s.rec = rec
		return true
	}
	return false
}

// Text returns the most recent record produced by Scan.
func (s *Scanner) Text() string { return s.rec }

// Err returns the first non-EOF error encountered by the underlying scanner.
func (s *Scanner) Err() error { return s.err }
