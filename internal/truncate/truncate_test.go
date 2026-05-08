package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/truncate"
)

func TestTruncateShortLine(t *testing.T) {
	tr := truncate.New(truncate.DefaultOptions())
	input := []byte("short line")
	out := tr.Truncate(input)
	if string(out) != "short line" {
		t.Fatalf("expected unchanged line, got %q", out)
	}
}

func TestTruncateExactLimit(t *testing.T) {
	opts := truncate.Options{MaxBytes: 10, AppendSuffix: false}
	tr := truncate.New(opts)
	input := []byte("1234567890")
	out := tr.Truncate(input)
	if string(out) != "1234567890" {
		t.Fatalf("expected exact match, got %q", out)
	}
}

func TestTruncateLongLineWithSuffix(t *testing.T) {
	opts := truncate.Options{MaxBytes: 10, AppendSuffix: true}
	tr := truncate.New(opts)
	input := []byte("hello world this is a long line")
	out := tr.Truncate(input)
	if len(out) > 10 {
		t.Fatalf("expected len <= 10, got %d: %q", len(out), out)
	}
	if !strings.HasSuffix(string(out), truncate.Suffix) {
		t.Fatalf("expected suffix %q, got %q", truncate.Suffix, out)
	}
}

func TestTruncateLongLineNoSuffix(t *testing.T) {
	opts := truncate.Options{MaxBytes: 8, AppendSuffix: false}
	tr := truncate.New(opts)
	input := []byte("abcdefghijklmnop")
	out := tr.Truncate(input)
	if len(out) != 8 {
		t.Fatalf("expected len 8, got %d", len(out))
	}
	if string(out) != "abcdefgh" {
		t.Fatalf("expected %q, got %q", "abcdefgh", out)
	}
}

func TestTruncateUTF8Boundary(t *testing.T) {
	// 'é' is 2 bytes (0xC3 0xA9); place it at the boundary
	opts := truncate.Options{MaxBytes: 6, AppendSuffix: false}
	tr := truncate.New(opts)
	// "abcdé" = 6 bytes; fits exactly
	input := []byte("abcd\xC3\xA9")
	out := tr.Truncate(input)
	if string(out) != "abcdé" {
		t.Fatalf("expected %q, got %q", "abcdé", out)
	}

	// "abcé" truncated at 5 bytes must not split the rune
	opts2 := truncate.Options{MaxBytes: 5, AppendSuffix: false}
	tr2 := truncate.New(opts2)
	out2 := tr2.Truncate(input)
	if !strings.HasPrefix(string(out2), "abc") {
		t.Fatalf("expected valid UTF-8 prefix, got %q", out2)
	}
}

func TestTruncateString(t *testing.T) {
	opts := truncate.Options{MaxBytes: 10, AppendSuffix: true}
	tr := truncate.New(opts)
	out := tr.TruncateString("hello world!!")
	if len(out) > 10 {
		t.Fatalf("expected len <= 10, got %d: %q", len(out), out)
	}
}

func TestTruncateZeroMaxBytes(t *testing.T) {
	// MaxBytes <= 0 should fall back to DefaultMaxBytes
	opts := truncate.Options{MaxBytes: 0, AppendSuffix: false}
	tr := truncate.New(opts)
	long := strings.Repeat("x", truncate.DefaultMaxBytes+100)
	out := tr.Truncate([]byte(long))
	if len(out) > truncate.DefaultMaxBytes {
		t.Fatalf("expected len <= %d, got %d", truncate.DefaultMaxBytes, len(out))
	}
}
