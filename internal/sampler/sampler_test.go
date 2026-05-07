package sampler

import (
	"testing"
)

func TestSamplerEveryOne(t *testing.T) {
	s, err := New(Options{Every: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 5; i++ {
		if !s.ShouldEmit() {
			t.Errorf("line %d: expected emit", i+1)
		}
	}
	if s.Emitted() != 5 {
		t.Errorf("expected 5 emitted, got %d", s.Emitted())
	}
}

func TestSamplerEveryThree(t *testing.T) {
	s, err := New(Options{Every: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, want := range expected {
		got := s.ShouldEmit()
		if got != want {
			t.Errorf("line %d: want %v got %v", i+1, want, got)
		}
	}
	if s.Emitted() != 3 {
		t.Errorf("expected 3 emitted, got %d", s.Emitted())
	}
}

func TestSamplerMaxLines(t *testing.T) {
	s, err := New(Options{Every: 1, MaxLines: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	emitted := 0
	for i := 0; i < 10; i++ {
		if s.ShouldEmit() {
			emitted++
		}
	}
	if emitted != 3 {
		t.Errorf("expected 3 emitted due to MaxLines, got %d", emitted)
	}
}

func TestSamplerReset(t *testing.T) {
	s, err := New(Options{Every: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s.ShouldEmit()
	s.ShouldEmit()
	s.Reset()
	if s.Seen() != 0 || s.Emitted() != 0 {
		t.Errorf("expected zeroed counters after reset, got seen=%d emitted=%d", s.Seen(), s.Emitted())
	}
}

func TestSamplerInvalidEvery(t *testing.T) {
	_, err := New(Options{Every: 0})
	if err == nil {
		t.Error("expected error for Every=0, got nil")
	}
}

func TestSamplerSeenCounter(t *testing.T) {
	s, _ := New(Options{Every: 5})
	for i := 0; i < 7; i++ {
		s.ShouldEmit()
	}
	if s.Seen() != 7 {
		t.Errorf("expected Seen()=7, got %d", s.Seen())
	}
	if s.Emitted() != 1 {
		t.Errorf("expected Emitted()=1, got %d", s.Emitted())
	}
}
