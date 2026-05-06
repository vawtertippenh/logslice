package highlight_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestApplyNoPattern(t *testing.T) {
	h := highlight.New(nil, highlight.Red)
	line := "2024-01-01 no match here"
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApplyNoMatch(t *testing.T) {
	h := highlight.New(regexp.MustCompile(`ERROR`), highlight.Red)
	line := "2024-01-01 INFO everything is fine"
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApplySingleMatch(t *testing.T) {
	h := highlight.New(regexp.MustCompile(`ERROR`), highlight.Red)
	line := "2024-01-01 ERROR something failed"
	got := h.Apply(line)
	if !strings.Contains(got, highlight.Red+"ERROR"+highlight.Reset) {
		t.Errorf("expected colored ERROR in output, got %q", got)
	}
	// Original text outside match must be preserved.
	if !strings.Contains(got, "2024-01-01 ") {
		t.Errorf("prefix text lost in output: %q", got)
	}
	if !strings.Contains(got, " something failed") {
		t.Errorf("suffix text lost in output: %q", got)
	}
}

func TestApplyMultipleMatches(t *testing.T) {
	h := highlight.New(regexp.MustCompile(`WARN`), highlight.Yellow)
	line := "WARN first WARN second WARN third"
	got := h.Apply(line)
	count := strings.Count(got, highlight.Yellow+"WARN"+highlight.Reset)
	if count != 3 {
		t.Errorf("expected 3 highlighted matches, got %d in %q", count, got)
	}
}

func TestApplyPreservesNonMatchContent(t *testing.T) {
	h := highlight.New(regexp.MustCompile(`\d+`), highlight.Cyan)
	line := "abc 123 def 456 ghi"
	got := h.Apply(line)
	stripped := strings.ReplaceAll(got, highlight.Cyan, "")
	stripped = strings.ReplaceAll(stripped, highlight.Reset, "")
	if stripped != line {
		t.Errorf("non-ANSI content changed: got %q, want %q", stripped, line)
	}
}

func TestEnabled(t *testing.T) {
	if highlight.New(nil, highlight.Red).Enabled() {
		t.Error("expected Enabled=false for nil pattern")
	}
	if !highlight.New(regexp.MustCompile(`x`), highlight.Red).Enabled() {
		t.Error("expected Enabled=true for non-nil pattern")
	}
}
