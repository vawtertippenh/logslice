package filter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/logline"
)

func makeTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestFilterMatchTimeRange(t *testing.T) {
	from := makeTime("2024-01-01T10:00:00Z")
	to := makeTime("2024-01-01T12:00:00Z")
	f := filter.New(filter.Options{From: from, To: to})

	inside := logline.LogLine{Timestamp: makeTime("2024-01-01T11:00:00Z"), Raw: "inside"}
	outside := logline.LogLine{Timestamp: makeTime("2024-01-01T13:00:00Z"), Raw: "outside"}

	if !f.Match(inside) {
		t.Error("expected inside line to match")
	}
	if f.Match(outside) {
		t.Error("expected outside line not to match")
	}
}

func TestFilterMatchRegex(t *testing.T) {
	pat := regexp.MustCompile(`ERROR`)
	f := filter.New(filter.Options{Pattern: pat})

	errLine := logline.LogLine{Raw: "2024-01-01 ERROR something failed"}
	infoLine := logline.LogLine{Raw: "2024-01-01 INFO all good"}

	if !f.Match(errLine) {
		t.Error("expected error line to match")
	}
	if f.Match(infoLine) {
		t.Error("expected info line not to match")
	}
}

func TestFilterMatchBoth(t *testing.T) {
	from := makeTime("2024-01-01T10:00:00Z")
	to := makeTime("2024-01-01T12:00:00Z")
	pat := regexp.MustCompile(`WARN`)
	f := filter.New(filter.Options{From: from, To: to, Pattern: pat})

	match := logline.LogLine{Timestamp: makeTime("2024-01-01T11:00:00Z"), Raw: "WARN disk low"}
	wrongTime := logline.LogLine{Timestamp: makeTime("2024-01-01T09:00:00Z"), Raw: "WARN disk low"}
	wrongPat := logline.LogLine{Timestamp: makeTime("2024-01-01T11:00:00Z"), Raw: "INFO all good"}

	if !f.Match(match) {
		t.Error("expected line to match both criteria")
	}
	if f.Match(wrongTime) {
		t.Error("expected line with wrong time not to match")
	}
	if f.Match(wrongPat) {
		t.Error("expected line with wrong pattern not to match")
	}
}

func TestFilterMatchString(t *testing.T) {
	pat := regexp.MustCompile(`DEBUG`)
	f := filter.New(filter.Options{Pattern: pat})

	if !f.MatchString("DEBUG verbose output") {
		t.Error("expected debug string to match")
	}
	if f.MatchString("INFO normal") {
		t.Error("expected info string not to match")
	}
}
