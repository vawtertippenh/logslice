package scanner_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/scanner"
)

// sampleLogs contains representative log lines for testing.
const sampleLogs = `2024-01-15 10:00:01 INFO  starting application
2024-01-15 10:00:02 DEBUG initializing database connection
2024-01-15 10:05:30 INFO  server listening on :8080
2024-01-15 10:10:00 WARN  high memory usage detected
2024-01-15 10:15:45 ERROR failed to process request: timeout
2024-01-15 10:20:00 INFO  processed 1000 requests
2024-01-15 10:25:10 DEBUG cache hit ratio: 0.85
2024-01-15 10:30:00 INFO  shutting down gracefully
`

func makeTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		t.Fatalf("makeTime: failed to parse %q: %v", value, err)
	}
	return parsed
}

func TestScannerNew(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	s := scanner.New(r, nil)
	if s == nil {
		t.Fatal("expected non-nil scanner")
	}
}

func TestScannerScanAll(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	s := scanner.New(r, nil)

	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	const wantCount = 8
	if len(lines) != wantCount {
		t.Errorf("got %d lines, want %d", len(lines), wantCount)
	}
}

func TestScannerTimeRangeFilter(t *testing.T) {
	r := strings.NewReader(sampleLogs)

	start := makeTime(t, "2024-01-15 10:05:00")
	end := makeTime(t, "2024-01-15 10:20:00")

	opts := &scanner.Options{
		Start: start,
		End:   end,
	}
	s := scanner.New(r, opts)

	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	// Expect lines at 10:05:30, 10:10:00, 10:15:45, 10:20:00
	const wantCount = 4
	if len(lines) != wantCount {
		t.Errorf("got %d lines, want %d; lines: %v", len(lines), wantCount, lines)
	}
}

func TestScannerEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	s := scanner.New(r, nil)

	var count int
	for s.Scan() {
		count++
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 lines from empty input, got %d", count)
	}
}

func TestScannerLargeInput(t *testing.T) {
	// Build a large log file with 10000 lines.
	var buf bytes.Buffer
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10000; i++ {
		ts := base.Add(time.Duration(i) * time.Second)
		buf.WriteString(ts.Format("2006-01-02 15:04:05") + " INFO line\n")
	}

	s := scanner.New(&buf, nil)
	var count int
	for s.Scan() {
		count++
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 10000 {
		t.Errorf("expected 10000 lines, got %d", count)
	}
}

func TestScannerNoTimestampLines(t *testing.T) {
	// Lines without timestamps should still be returned when no filter is set.
	input := "no timestamp here\nanother plain line\n"
	r := strings.NewReader(input)
	s := scanner.New(r, nil)

	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
