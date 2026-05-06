package rotate

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatcherDetectsRotation(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/app.log"
	if err := os.WriteFile(path, []byte("initial\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := NewWatcher(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go w.Run(ctx)

	// Truncate the file to trigger rotation detection.
	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev, ok := <-w.Events:
		if !ok {
			t.Fatal("events channel closed unexpectedly")
		}
		if ev != EventRotated {
			t.Errorf("expected EventRotated, got %v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for rotation event")
	}
}

func TestWatcherNoEventWhenStable(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/stable.log"
	if err := os.WriteFile(path, []byte("stable\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := NewWatcher(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	go w.Run(ctx)

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event %v on stable file", ev)
	case <-ctx.Done():
		// expected: no events
	}
}
