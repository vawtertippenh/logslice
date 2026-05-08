// Package multiline provides support for aggregating multi-line log entries
// into a single logical record.
//
// Many log formats (Java stack traces, Python tracebacks, structured JSON blobs
// split across lines) emit a single event across several physical lines.  The
// Aggregator in this package detects continuation lines via a user-supplied
// regular expression and buffers them until the next "start" line (or EOF)
// before emitting the completed record.
//
// Usage:
//
//	agg, err := multiline.New(multiline.Options{
//		StartPattern: `^\d{4}-\d{2}-\d{2}`,
//		MaxLines:     500,
//	})
//	if err != nil { ... }
//
//	for _, raw := range lines {
//		if record, ok := agg.Feed(raw); ok {
//			// record is a complete, joined log entry
//		}
//	}
//	if record, ok := agg.Flush(); ok {
//		// last buffered record
//	}
package multiline
