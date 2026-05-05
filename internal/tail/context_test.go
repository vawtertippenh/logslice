package tail

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/filter"
)

func makeFilter(t *testing.T, pattern string) *filter.Filter {
	t.Helper()
	f, err := filter.New(filter.Options{Regex: pattern})
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}
	return f
}

func TestContextReaderMatch(t *testing.T) {
	input := "aaa\nbbb\nccc ERROR foo\nddd\neee\n"
	r := New(Options{NumLines: 20, ChunkSize: 64})
	f := makeFilter(t, "ERROR")
	cr := NewContextReader(r, f, 1, 1)

	results, err := cr.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("want 1 result, got %d", len(results))
	}
	res := results[0]
	if res.Line != "ccc ERROR foo" {
		t.Errorf("unexpected match line: %q", res.Line)
	}
	if len(res.Before) != 1 || res.Before[0] != "bbb" {
		t.Errorf("unexpected before: %v", res.Before)
	}
	if len(res.After) != 1 || res.After[0] != "ddd" {
		t.Errorf("unexpected after: %v", res.After)
	}
}

func TestContextReaderNoMatch(t *testing.T) {
	input := "line1\nline2\nline3\n"
	r := New(Options{NumLines: 10})
	f := makeFilter(t, "NOTFOUND")
	cr := NewContextReader(r, f, 2, 2)

	results, err := cr.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestContextResultFormat(t *testing.T) {
	cr := ContextResult{
		Before: []string{"prev"},
		Line:   "match",
		After:  []string{"next"},
	}
	out := cr.Format()
	if !strings.Contains(out, "» match") {
		t.Errorf("format missing match indicator: %q", out)
	}
	if !strings.Contains(out, "  prev") {
		t.Errorf("format missing before line: %q", out)
	}
	if !strings.Contains(out, "  next") {
		t.Errorf("format missing after line: %q", out)
	}
}
