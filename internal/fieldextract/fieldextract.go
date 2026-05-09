// Package fieldextract provides utilities for extracting named fields
// from structured log lines (e.g. key=value or key="value" pairs).
package fieldextract

import (
	"regexp"
	"strings"
)

// Field represents a single extracted key-value pair from a log line.
type Field struct {
	Key   string
	Value string
}

// Extractor parses structured fields from log lines.
type Extractor struct {
	pairRe *regexp.Regexp
}

// New returns an Extractor ready to parse key=value pairs.
// It handles both bare values (key=val) and quoted values (key="val with spaces").
func New() *Extractor {
	// Matches: key="quoted value" or key=bare_value
	re := regexp.MustCompile(`(\w[\w.\-]*)=("[^"]*"|\S+)`)
	return &Extractor{pairRe: re}
}

// Extract returns all key=value fields found in line.
func (e *Extractor) Extract(line string) []Field {
	matches := e.pairRe.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return nil
	}
	fields := make([]Field, 0, len(matches))
	for _, m := range matches {
		val := m[2]
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		fields = append(fields, Field{Key: m[1], Value: val})
	}
	return fields
}

// ExtractMap returns extracted fields as a map. If a key appears multiple
// times the last occurrence wins.
func (e *Extractor) ExtractMap(line string) map[string]string {
	fields := e.Extract(line)
	if len(fields) == 0 {
		return nil
	}
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f.Key] = f.Value
	}
	return m
}

// Get extracts the value for a single named key from line.
// Returns the value and true if found, or empty string and false otherwise.
func (e *Extractor) Get(line, key string) (string, bool) {
	for _, f := range e.Extract(line) {
		if strings.EqualFold(f.Key, key) {
			return f.Value, true
		}
	}
	return "", false
}
