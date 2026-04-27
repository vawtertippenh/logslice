package filter

import (
	"regexp"
	"time"

	"github.com/yourorg/logslice/internal/logline"
)

// Options holds the filtering criteria for log lines.
type Options struct {
	From    time.Time
	To      time.Time
	Pattern *regexp.Regexp
}

// Filter applies time-range and regex criteria to log lines.
type Filter struct {
	opts Options
}

// New creates a new Filter with the given options.
func New(opts Options) *Filter {
	return &Filter{opts: opts}
}

// Match returns true if the log line satisfies all active filter criteria.
func (f *Filter) Match(line logline.LogLine) bool {
	if !f.opts.From.IsZero() || !f.opts.To.IsZero() {
		if !line.InRange(f.opts.From, f.opts.To) {
			return false
		}
	}
	if f.opts.Pattern != nil {
		if !f.opts.Pattern.MatchString(line.Raw) {
			return false
		}
	}
	return true
}

// MatchString returns true if the raw string satisfies the regex filter.
// It ignores time-range checks since no timestamp is available.
func (f *Filter) MatchString(raw string) bool {
	if f.opts.Pattern == nil {
		return true
	}
	return f.opts.Pattern.MatchString(raw)
}
