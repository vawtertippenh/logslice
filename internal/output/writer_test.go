package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestWriterWriteLine(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	lines := []string{
		"2024-01-01 00:00:01 INFO starting server",
		"2024-01-01 00:00:02 DEBUG connection accepted",
		"2024-01-01 00:00:03 ERROR disk full",
	}

	for _, l := range lines {
		if err := w.WriteLine([]byte(l)); err != nil {
			t.Fatalf("WriteLine(%q) unexpected error: %v", l, err)
		}
	}

	if err := w.Flush(); err != nil {
		t.Fatalf("Flush() unexpected error: %v", err)
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("output missing line %q", l)
		}
	}
}

func TestWriterLinesWritten(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	for i := 0; i < 5; i++ {
		if err := w.WriteLine([]byte("some log line")); err != nil {
			t.Fatalf("WriteLine error: %v", err)
		}
	}
	_ = w.Flush()

	if got := w.LinesWritten(); got != 5 {
		t.Errorf("LinesWritten() = %d, want 5", got)
	}
}

func TestWriterBytesWritten(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	line := []byte("hello world")
	_ = w.WriteLine(line)
	_ = w.Flush()

	// "hello world\n" = 12 bytes
	if got := w.BytesWritten(); got != 12 {
		t.Errorf("BytesWritten() = %d, want 12", got)
	}
}

func TestWriterStats(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)
	_ = w.WriteLine([]byte("test"))
	_ = w.Flush()

	stats := w.Stats()
	if !strings.Contains(stats, "lines=1") {
		t.Errorf("Stats() = %q, expected to contain 'lines=1'", stats)
	}
	if !strings.Contains(stats, "bytes=") {
		t.Errorf("Stats() = %q, expected to contain 'bytes='", stats)
	}
}
