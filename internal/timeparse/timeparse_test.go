package timeparse

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		wantYear int
		wantMonth time.Month
		wantDay  int
	}{
		{"2024-03-15T08:30:00Z", false, 2024, time.March, 15},
		{"2024-03-15T08:30:00.123+02:00", false, 2024, time.March, 15},
		{"2024-03-15 08:30:00", false, 2024, time.March, 15},
		{"2024-03-15 08:30:00.456", false, 2024, time.March, 15},
		{"2024/03/15 08:30:00", false, 2024, time.March, 15},
		{"15/Mar/2024:08:30:00 +0000", false, 2024, time.March, 15},
		{"not-a-timestamp", true, 0, 0, 0},
		{"", true, 0, 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, layout, err := Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil (layout=%s)", tc.input, layout)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}
			if got.Year() != tc.wantYear || got.Month() != tc.wantMonth || got.Day() != tc.wantDay {
				t.Errorf("Parse(%q) = %v, want %d-%02d-%02d",
					tc.input, got, tc.wantYear, tc.wantMonth, tc.wantDay)
			}
		})
	}
}

func TestParseWithLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	got, _, err := ParseWithLocation("2024-03-15 12:00:00", loc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, offset := got.Zone()
	if offset == 0 {
		t.Errorf("expected non-UTC offset for America/New_York, got 0")
	}
}

func TestKnownFormats(t *testing.T) {
	formats := KnownFormats()
	if len(formats) == 0 {
		t.Error("KnownFormats returned empty slice")
	}
	// Mutating the returned slice should not affect internals.
	formats[0] = "corrupted"
	if KnownFormats()[0] == "corrupted" {
		t.Error("KnownFormats returned a reference to internal slice")
	}
}
