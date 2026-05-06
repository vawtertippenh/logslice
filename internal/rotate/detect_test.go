package rotate

import (
	"os"
	"testing"
)

func TestDetectorNoRotation(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("hello\n")
	f.Close()

	det, err := NewDetector(f.Name())
	if err != nil {
		t.Fatalf("NewDetector: %v", err)
	}

	rotated, err := det.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if rotated {
		t.Error("expected not rotated, got rotated")
	}
}

func TestDetectorTruncation(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/test.log"

	if err := os.WriteFile(path, []byte("line1\nline2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	det, err := NewDetector(path)
	if err != nil {
		t.Fatal(err)
	}

	// Truncate the file.
	if err := os.WriteFile(path, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	rotated, err := det.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if !rotated {
		t.Error("expected rotated after truncation")
	}
}

func TestDetectorMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/test.log"
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	det, err := NewDetector(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	rotated, err := det.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Error("expected rotated when file is missing")
	}
}

func TestDetectorReset(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/test.log"
	if err := os.WriteFile(path, []byte("line1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	det, err := NewDetector(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := det.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	rotated, err := det.Rotated()
	if err != nil {
		t.Fatal(err)
	}
	if rotated {
		t.Error("expected not rotated after reset")
	}
}
