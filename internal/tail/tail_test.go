package tail

import (
	"strings"
	"testing"
)

func makeReader(rs strings.Reader, n int) *Reader {
	return New(Options{NumLines: n, ChunkSize: 16})
}

func TestReadExactLines(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"
	r := New(Options{NumLines: 3, ChunkSize: 16})
	lines, err := r.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "line3" || lines[1] != "line4" || lines[2] != "line5" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestReadFewerLinesThanRequested(t *testing.T) {
	input := "only\ntwo\n"
	r := New(Options{NumLines: 10, ChunkSize: 64})
	lines, err := r.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
	}
}

func TestReadEmptyInput(t *testing.T) {
	r := New(Options{NumLines: 5})
	lines, err := r.Read(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected no lines, got %v", lines)
	}
}

func TestReadDefaultOptions(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 20; i++ {
		sb.WriteString("logline\n")
	}
	r := New(Options{})
	lines, err := r.Read(strings.NewReader(sb.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 10 {
		t.Errorf("want 10 lines (default), got %d", len(lines))
	}
}

func TestReadSmallChunk(t *testing.T) {
	input := "alpha\nbeta\ngamma\ndelta\n"
	r := New(Options{NumLines: 2, ChunkSize: 4})
	lines, err := r.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "gamma" || lines[1] != "delta" {
		t.Errorf("unexpected lines: %v", lines)
	}
}
