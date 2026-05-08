// Package truncate provides line truncation utilities for log output.
//
// It supports truncating long lines to a maximum byte or rune length,
// appending a configurable suffix (e.g. "...") to indicate truncation,
// and respects UTF-8 character boundaries to avoid producing invalid output.
//
// Example usage:
//
//	t := truncate.New(truncate.Options{Limit: 120, Suffix: "…"})
//	short, wasTruncated := t.Line(longLine)
package truncate
