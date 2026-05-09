package fieldextract_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldextract"
)

func TestExtractEmpty(t *testing.T) {
	e := fieldextract.New()
	fields := e.Extract("no fields here at all")
	if len(fields) != 0 {
		t.Fatalf("expected 0 fields, got %d", len(fields))
	}
}

func TestExtractBareValues(t *testing.T) {
	e := fieldextract.New()
	line := `level=info msg=started pid=1234`
	fields := e.Extract(line)
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	expect := []fieldextract.Field{
		{Key: "level", Value: "info"},
		{Key: "msg", Value: "started"},
		{Key: "pid", Value: "1234"},
	}
	for i, f := range fields {
		if f != expect[i] {
			t.Errorf("field[%d]: got %+v, want %+v", i, f, expect[i])
		}
	}
}

func TestExtractQuotedValue(t *testing.T) {
	e := fieldextract.New()
	line := `level=warn msg="disk usage high" host=web-01`
	fields := e.Extract(line)
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[1].Value != "disk usage high" {
		t.Errorf("expected quoted value stripped, got %q", fields[1].Value)
	}
}

func TestExtractMap(t *testing.T) {
	e := fieldextract.New()
	line := `level=debug component=scanner latency=42ms`
	m := e.ExtractMap(line)
	if m == nil {
		t.Fatal("expected non-nil map")
	}
	if m["level"] != "debug" {
		t.Errorf("level: got %q", m["level"])
	}
	if m["latency"] != "42ms" {
		t.Errorf("latency: got %q", m["latency"])
	}
}

func TestGet(t *testing.T) {
	e := fieldextract.New()
	line := `ts=2024-01-02T15:04:05Z level=error msg="oops"`
	v, ok := e.Get(line, "level")
	if !ok {
		t.Fatal("expected key 'level' to be found")
	}
	if v != "error" {
		t.Errorf("expected 'error', got %q", v)
	}
}

func TestGetMissing(t *testing.T) {
	e := fieldextract.New()
	_, ok := e.Get("level=info", "host")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestExtractMapEmpty(t *testing.T) {
	e := fieldextract.New()
	m := e.ExtractMap("plain log line with no pairs")
	if m != nil {
		t.Errorf("expected nil map for no-match input, got %v", m)
	}
}
