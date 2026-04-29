package index

import (
	"strings"
	"testing"
)

const sampleLog = `2024-01-01T10:00:00Z INFO starting server
2024-01-01T10:01:00Z INFO listening on :8080
2024-01-01T10:02:00Z DEBUG accepted connection
2024-01-01T10:03:00Z INFO request received
2024-01-01T10:04:00Z WARN slow response
2024-01-01T10:05:00Z ERROR timeout
`

func TestBuilderBuildSampleRate1(t *testing.T) {
	r := strings.NewReader(sampleLog)
	rs := &readSeeker{Reader: strings.NewReader(sampleLog), data: sampleLog}
	_ = r

	b := NewBuilder(1)
	idx, err := b.Build(rs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len() == 0 {
		t.Fatal("expected non-empty index")
	}
}

func TestBuilderBuildSampleRate2(t *testing.T) {
	rs := &readSeeker{data: sampleLog}

	b := NewBuilder(2)
	idx, err := b.Build(rs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With sample rate 2 we expect roughly half the entries
	if idx.Len() == 0 {
		t.Fatal("expected non-empty index")
	}
}

func TestBuilderEmptyInput(t *testing.T) {
	rs := &readSeeker{data: ""}
	b := NewBuilder(1)
	idx, err := b.Build(rs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len() != 0 {
		t.Errorf("expected empty index, got %d entries", idx.Len())
	}
}

// readSeeker wraps a string to implement io.ReadSeeker.
type readSeeker struct {
	*strings.Reader
	data string
}

func (rs *readSeeker) Read(p []byte) (int, error) {
	if rs.Reader == nil {
		rs.Reader = strings.NewReader(rs.data)
	}
	return rs.Reader.Read(p)
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	if rs.Reader == nil {
		rs.Reader = strings.NewReader(rs.data)
	}
	return rs.Reader.Seek(offset, whence)
}
