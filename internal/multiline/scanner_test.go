package multiline

import (
	"strings"
	"testing"
)

func TestScannerBasic(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01 event one",
		"  trace line 1",
		"  trace line 2",
		"2024-01-02 event two",
		"  trace line 3",
	}, "\n")

	s, err := NewScanner(strings.NewReader(input), Options{
		StartPattern: `^\d{4}-`,
	})
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}

	var records []string
	for s.Scan() {
		records = append(records, s.Text())
	}
	if s.Err() != nil {
		t.Fatalf("scanner error: %v", s.Err())
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d: %v", len(records), records)
	}
	if !strings.Contains(records[0], "trace line 2") {
		t.Errorf("record[0] missing continuation: %q", records[0])
	}
	if !strings.Contains(records[1], "trace line 3") {
		t.Errorf("record[1] missing continuation: %q", records[1])
	}
}

func TestScannerEmptyInput(t *testing.T) {
	s, _ := NewScanner(strings.NewReader(""), Options{StartPattern: `^X`})
	if s.Scan() {
		t.Fatal("expected no records for empty input")
	}
}

func TestScannerSingleLine(t *testing.T) {
	s, _ := NewScanner(strings.NewReader("START only"), Options{StartPattern: `^START`})
	if !s.Scan() {
		t.Fatal("expected one record")
	}
	if s.Text() != "START only" {
		t.Errorf("unexpected record: %q", s.Text())
	}
	if s.Scan() {
		t.Fatal("expected no more records")
	}
}

func TestScannerBadPattern(t *testing.T) {
	_, err := NewScanner(strings.NewReader("x"), Options{StartPattern: "["})
	if err == nil {
		t.Fatal("expected error for bad pattern")
	}
}
