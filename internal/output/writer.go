package output

import (
	"bufio"
	"fmt"
	"io"
	"sync/atomic"
)

// Writer wraps an io.Writer with buffering and line counting.
type Writer struct {
	bw        *bufio.Writer
	linesWritten atomic.Int64
	bytesWritten atomic.Int64
}

// New creates a new Writer wrapping the given io.Writer.
func New(w io.Writer) *Writer {
	return &Writer{
		bw: bufio.NewWriterSize(w, 64*1024),
	}
}

// WriteLine writes a single log line followed by a newline character.
func (w *Writer) WriteLine(line []byte) error {
	n, err := fmt.Fprintf(w.bw, "%s\n", line)
	if err != nil {
		return fmt.Errorf("output: write line: %w", err)
	}
	w.linesWritten.Add(1)
	w.bytesWritten.Add(int64(n))
	return nil
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	if err := w.bw.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	return nil
}

// LinesWritten returns the total number of lines written.
func (w *Writer) LinesWritten() int64 {
	return w.linesWritten.Load()
}

// BytesWritten returns the total number of bytes written.
func (w *Writer) BytesWritten() int64 {
	return w.bytesWritten.Load()
}

// Stats returns a human-readable summary of written output.
func (w *Writer) Stats() string {
	return fmt.Sprintf("lines=%d bytes=%d", w.LinesWritten(), w.BytesWritten())
}
