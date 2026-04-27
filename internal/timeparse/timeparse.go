package timeparse

import (
	"fmt"
	"time"
)

// Common log timestamp formats to try when parsing.
var knownFormats = []string{
	"2006-01-02T15:04:05Z07:00",   // RFC3339
	"2006-01-02T15:04:05.000Z07:00", // RFC3339 with millis
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700", // Common HTTP access log format
	"Jan 02 15:04:05",            // syslog
	"Jan  2 15:04:05",            // syslog (single-digit day)
}

// Parse attempts to parse a timestamp string using a list of known formats.
// It returns the parsed time and the matched format string, or an error if
// none of the known formats match.
func Parse(s string) (time.Time, string, error) {
	for _, layout := range knownFormats {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparse: unrecognized timestamp format: %q", s)
}

// ParseWithLocation is like Parse but applies the given location when the
// format does not include timezone information.
func ParseWithLocation(s string, loc *time.Location) (time.Time, string, error) {
	for _, layout := range knownFormats {
		t, err := time.ParseInLocation(layout, s, loc)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparse: unrecognized timestamp format: %q", s)
}

// KnownFormats returns a copy of the list of supported timestamp formats.
func KnownFormats() []string {
	out := make([]string, len(knownFormats))
	copy(out, knownFormats)
	return out
}
