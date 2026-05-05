package tail

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ContextResult holds a matched line along with surrounding context lines.
type ContextResult struct {
	// Before contains lines preceding the match.
	Before []string
	// Line is the matched line.
	Line string
	// After contains lines following the match.
	After []string
}

// ContextReader reads lines from the tail and returns matches with context.
type ContextReader struct {
	reader  *Reader
	f       *filter.Filter
	before  int
	after   int
}

// NewContextReader creates a ContextReader that applies the given filter and
// attaches up to before/after lines of context around each match.
func NewContextReader(r *Reader, f *filter.Filter, before, after int) *ContextReader {
	return &ContextReader{reader: r, f: f, before: before, after: after}
}

// Read reads from rs and returns all matching lines with context.
func (c *ContextReader) Read(rs io.ReadSeeker) ([]ContextResult, error) {
	lines, err := c.reader.Read(rs)
	if err != nil {
		return nil, fmt.Errorf("context: read: %w", err)
	}
	var results []ContextResult
	for i, line := range lines {
		if !c.f.MatchString(line) {
			continue
		}
		cr := ContextResult{Line: line}
		start := i - c.before
		if start < 0 {
			start = 0
		}
		end := i + c.after + 1
		if end > len(lines) {
			end = len(lines)
		}
		if i > 0 && start < i {
			cr.Before = append([]string(nil), lines[start:i]...)
		}
		if i+1 < end {
			cr.After = append([]string(nil), lines[i+1:end]...)
		}
		results = append(results, cr)
	}
	return results, nil
}

// Format returns a human-readable representation of a ContextResult.
func (cr ContextResult) Format() string {
	var sb strings.Builder
	for _, b := range cr.Before {
		sb.WriteString("  " + b + "\n")
	}
	sb.WriteString("» " + cr.Line + "\n")
	for _, a := range cr.After {
		sb.WriteString("  " + a + "\n")
	}
	return sb.String()
}
