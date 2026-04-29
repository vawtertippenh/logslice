package index

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func makeSearcherIndex() *Index {
	idx := New()
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		idx.Add(base.Add(time.Duration(i)*time.Minute), int64(i*100))
	}
	return idx
}

func TestSearcherSeekToStart(t *testing.T) {
	idx := makeSearcherIndex()
	s := NewSearcher(idx)

	data := make([]byte, 500)
	rs := bytes.NewReader(data)

	start := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)
	offset, err := s.SeekToStart(rs, start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// FindStart returns offset of entry at or before 10:02, which is entry index 2 => offset 200
	if offset != 200 {
		t.Errorf("expected offset 200, got %d", offset)
	}
	// Confirm the reader position matches
	pos, _ := rs.Seek(0, io.SeekCurrent)
	if pos != offset {
		t.Errorf("reader position %d does not match offset %d", pos, offset)
	}
}

func TestSearcherSeekToStartBeforeAll(t *testing.T) {
	idx := makeSearcherIndex()
	s := NewSearcher(idx)

	data := make([]byte, 500)
	rs := bytes.NewReader(data)

	start := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	offset, err := s.SeekToStart(rs, start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0 for time before all entries, got %d", offset)
	}
}

func TestSearcherSeekToEnd(t *testing.T) {
	idx := makeSearcherIndex()
	s := NewSearcher(idx)

	data := make([]byte, 500)
	rs := bytes.NewReader(data)

	end := time.Date(2024, 1, 1, 10, 2, 30, 0, time.UTC)
	endOffset := s.SeekToEnd(rs, end)
	// First entry strictly after 10:02:30 is 10:03 => offset 300
	if endOffset != 300 {
		t.Errorf("expected end offset 300, got %d", endOffset)
	}
}

func TestSearcherSeekToStartAndEnd(t *testing.T) {
	idx := makeSearcherIndex()
	s := NewSearcher(idx)

	data := make([]byte, 500)
	rs := bytes.NewReader(data)

	start := time.Date(2024, 1, 1, 10, 1, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 10, 3, 0, 0, time.UTC)

	startOff, endOff, err := s.SeekToStartAndEnd(rs, start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startOff != 100 {
		t.Errorf("expected startOff 100, got %d", startOff)
	}
	if endOff != 400 {
		t.Errorf("expected endOff 400, got %d", endOff)
	}
}
