package logline

import (
	"time"
)

// LogLine represents a single parsed log line with its timestamp and raw content.
type LogLine struct {
	Timestamp time.Time
	Raw       string
	Offset    int64
}

// IsZero reports whether the log line has no valid timestamp.
func (l LogLine) IsZero() bool {
	return l.Timestamp.IsZero()
}

// InRange reports whether the log line's timestamp falls within [start, end].
// If start or end is zero, that bound is considered open.
func (l LogLine) InRange(start, end time.Time) bool {
	if l.IsZero() {
		return false
	}
	if !start.IsZero() && l.Timestamp.Before(start) {
		return false
	}
	if !end.IsZero() && l.Timestamp.After(end) {
		return false
	}
	return true
}
