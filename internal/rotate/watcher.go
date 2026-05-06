package rotate

import (
	"context"
	"time"
)

// Event describes what kind of rotation was observed.
type Event int

const (
	EventNone     Event = iota
	EventRotated        // file replaced or inode changed
	EventTruncated      // same inode but size shrank
)

// Watcher polls a file for rotation events and emits them on a channel.
type Watcher struct {
	detector *Detector
	interval time.Duration
	Events   <-chan Event
	events   chan Event
}

// NewWatcher creates a Watcher that polls path every interval.
func NewWatcher(path string, interval time.Duration) (*Watcher, error) {
	det, err := NewDetector(path)
	if err != nil {
		return nil, err
	}
	ch := make(chan Event, 4)
	return &Watcher{
		detector: det,
		interval: interval,
		Events:   ch,
		events:   ch,
	}, nil
}

// Run starts polling until ctx is cancelled. It is safe to call in a goroutine.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	defer close(w.events)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rotated, err := w.detector.Rotated()
			if err != nil || !rotated {
				continue
			}
			ev := EventRotated
			_ = w.detector.Reset()
			select {
			case w.events <- ev:
			default:
			}
		}
	}
}
