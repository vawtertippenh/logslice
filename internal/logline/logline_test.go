package logline

import (
	"testing"
	"time"
)

func TestLogLineIsZero(t *testing.T) {
	ll := LogLine{Raw: "some log line"}
	if !ll.IsZero() {
		t.Error("expected IsZero() == true for line without timestamp")
	}

	ll.Timestamp = time.Now()
	if ll.IsZero() {
		t.Error("expected IsZero() == false for line with timestamp")
	}
}

func TestLogLineInRange(t *testing.T) {
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	ll := LogLine{Timestamp: base, Raw: "msg"}

	start := base.Add(-time.Hour)
	end := base.Add(time.Hour)

	if !ll.InRange(start, end) {
		t.Error("expected InRange == true")
	}
	if ll.InRange(base.Add(time.Minute), end) {
		t.Error("expected InRange == false when timestamp before start")
	}
	if ll.InRange(start, base.Add(-time.Minute)) {
		t.Error("expected InRange == false when timestamp after end")
	}

	// Open bounds
	if !ll.InRange(time.Time{}, end) {
		t.Error("expected InRange == true with open start")
	}
	if !ll.InRange(start, time.Time{}) {
		t.Error("expected InRange == true with open end")
	}

	// Zero timestamp line
	zero := LogLine{Raw: "no ts"}
	if zero.InRange(time.Time{}, time.Time{}) {
		t.Error("expected InRange == false for zero-timestamp line")
	}
}

func TestLogLineInRangeBoundaryExact(t *testing.T) {
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	ll := LogLine{Timestamp: base, Raw: "msg"}

	// Timestamp exactly equal to start should be in range
	if !ll.InRange(base, base.Add(time.Hour)) {
		t.Error("expected InRange == true when timestamp equals start")
	}

	// Timestamp exactly equal to end should be in range
	if !ll.InRange(base.Add(-time.Hour), base) {
		t.Error("expected InRange == true when timestamp equals end")
	}
}
