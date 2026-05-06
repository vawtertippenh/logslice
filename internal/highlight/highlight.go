// Package highlight provides ANSI terminal colorization for matched
// substrings in log lines, making regex matches visually distinct in
// interactive output sessions.
package highlight

import (
	"regexp"
	"strings"
)

const (
	// ANSI escape codes for colorization.
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Highlighter wraps matched substrings with ANSI color codes.
type Highlighter struct {
	pattern *regexp.Regexp
	prefix  string
	suffix  string
}

// New returns a Highlighter that marks matches of pattern with the
// given ANSI prefix and the reset suffix. If pattern is nil every
// call to Apply returns the line unchanged.
func New(pattern *regexp.Regexp, colorCode string) *Highlighter {
	return &Highlighter{
		pattern: pattern,
		prefix:  colorCode,
		suffix:  Reset,
	}
}

// Apply returns line with every match of the underlying pattern
// wrapped in the configured ANSI color codes. If no pattern was
// provided, or if the line contains no matches, the original string
// is returned without allocation.
func (h *Highlighter) Apply(line string) string {
	if h.pattern == nil {
		return line
	}
	loc := h.pattern.FindStringIndex(line)
	if loc == nil {
		return line
	}
	var sb strings.Builder
	sb.Grow(len(line) + 16)
	last := 0
	for _, idx := range h.pattern.FindAllStringIndex(line, -1) {
		sb.WriteString(line[last:idx[0]])
		sb.WriteString(h.prefix)
		sb.WriteString(line[idx[0]:idx[1]])
		sb.WriteString(h.suffix)
		last = idx[1]
	}
	sb.WriteString(line[last:])
	return sb.String()
}

// Enabled reports whether the Highlighter will actually modify lines.
func (h *Highlighter) Enabled() bool {
	return h.pattern != nil
}

// Strip removes all ANSI escape sequences from s, returning plain text.
// This is useful when writing highlighted output to a file or a
// non-terminal destination where escape codes would appear as raw bytes.
func Strip(s string) string {
	var sb strings.Builder
	for len(s) > 0 {
		idx := strings.IndexByte(s, '\033')
		if idx == -1 {
			sb.WriteString(s)
			break
		}
		sb.WriteString(s[:idx])
		// Skip past the escape sequence: ESC '[' ... letter
		s = s[idx:]
		end := strings.IndexFunc(s, func(r rune) bool {
			return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
		})
		if end == -1 {
			break
		}
		s = s[end+1:]
	}
	return sb.String()
}
