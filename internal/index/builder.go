package index

import (
	"bufio"
	"io"
	"regexp"
	"time"

	"github.com/yourusername/logslice/internal/timeparse"
)

// Builder scans a log file and builds an Index by sampling every N lines.
type Builder struct {
	sampleRate int
	pattern    *regexp.Regexp
}

// NewBuilder creates a Builder that indexes every sampleRate-th line.
// sampleRate must be >= 1.
func NewBuilder(sampleRate int) *Builder {
	if sampleRate < 1 {
		sampleRate = 1
	}
	return &Builder{
		sampleRate: sampleRate,
		pattern:    regexp.MustCompile(`^(\S+\s+\S+|\S+T\S+)`),
	}
}

// Build reads from r and returns an Index with sampled offset entries.
func (b *Builder) Build(r io.ReadSeeker) (*Index, error) {
	idx := New()

	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var offset int64
	lineNum := 0
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		lineLen := int64(len(line)) + 1 // +1 for newline

		if lineNum%b.sampleRate == 0 {
			ts := b.extractTime(line)
			if !ts.IsZero() {
				idx.Add(ts, offset)
			}
		}

		offset += lineLen
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return idx, nil
}

// extractTime attempts to parse a timestamp from the beginning of a log line.
func (b *Builder) extractTime(line string) time.Time {
	match := b.pattern.FindString(line)
	if match == "" {
		return time.Time{}
	}
	ts, err := timeparse.Parse(match)
	if err != nil {
		return time.Time{}
	}
	return ts
}
