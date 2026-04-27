package logline

import (
	"testing"
	"time"
)

func TestParserParse(t *testing.T) {
	p := NewParser(time.UTC)

	tests := []struct {
		name      string
		raw       string
		wantZero  bool
	}{
		{
			name:     "RFC3339 prefix",
			raw:      "2024-01-15T12:00:00Z INFO server started",
			wantZero: false,
		},
		{
			name:     "no timestamp",
			raw:      "this line has no timestamp at all",
			wantZero: true,
		},
		{
			name:     "empty line",
			raw:      "",
			wantZero: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ll := p.Parse(tc.raw, 0)
			if ll.Raw != tc.raw {
				t.Errorf("Raw mismatch: got %q want %q", ll.Raw, tc.raw)
			}
			if tc.wantZero && !ll.IsZero() {
				t.Errorf("expected zero timestamp for line: %q", tc.raw)
			}
			if !tc.wantZero && ll.IsZero() {
				t.Errorf("expected non-zero timestamp for line: %q", tc.raw)
			}
		})
	}
}

func TestParserOffset(t *testing.T) {
	p := NewParser(nil)
	ll := p.Parse("2024-06-01T00:00:00Z test", 1024)
	if ll.Offset != 1024 {
		t.Errorf("expected offset 1024, got %d", ll.Offset)
	}
}
