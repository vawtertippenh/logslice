package logline

import (
	"strings"
	"time"

	"github.com/user/logslice/internal/timeparse"
)

// Parser extracts timestamps from raw log lines.
type Parser struct {
	location *time.Location
	// prefixLen is the number of characters to inspect at the start of each line.
	prefixLen int
}

// NewParser creates a Parser using the given location for timestamp parsing.
// If loc is nil, time.UTC is used.
func NewParser(loc *time.Location) *Parser {
	if loc == nil {
		loc = time.UTC
	}
	return &Parser{
		location:  loc,
		prefixLen: 40,
	}
}

// Parse attempts to extract a timestamp from the beginning of the raw line.
// If no timestamp is found, the returned LogLine has a zero Timestamp.
func (p *Parser) Parse(raw string, offset int64) LogLine {
	candidate := raw
	if len(candidate) > p.prefixLen {
		candidate = candidate[:p.prefixLen]
	}
	candidate = strings.TrimSpace(candidate)

	t, err := timeparse.ParseWithLocation(candidate, p.location)
	if err != nil {
		// Try progressively shorter substrings to find a timestamp prefix.
		for i := len(candidate) - 1; i > 0; i-- {
			t, err = timeparse.ParseWithLocation(candidate[:i], p.location)
			if err == nil {
				break
			}
		}
	}

	return LogLine{
		Timestamp: t,
		Raw:       raw,
		Offset:    offset,
	}
}
