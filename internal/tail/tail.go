// Package tail provides functionality to read the last N lines of a log file
// efficiently using backward seeking, without reading the entire file.
package tail

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const defaultChunkSize = 4096

// Options configures the tail reader.
type Options struct {
	// NumLines is the number of lines to read from the end of the file.
	NumLines int
	// ChunkSize is the size of each backward-read chunk in bytes.
	ChunkSize int64
}

// Reader reads the last N lines from a file.
type Reader struct {
	opts Options
}

// New creates a new tail Reader with the given options.
func New(opts Options) *Reader {
	if opts.ChunkSize <= 0 {
		opts.ChunkSize = defaultChunkSize
	}
	if opts.NumLines <= 0 {
		opts.NumLines = 10
	}
	return &Reader{opts: opts}
}

// ReadFile returns the last N lines from the named file.
func (r *Reader) ReadFile(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("tail: open %q: %w", name, err)
	}
	defer f.Close()
	return r.Read(f)
}

// Read returns the last N lines from the given ReadSeeker.
func (r *Reader) Read(rs io.ReadSeeker) ([]string, error) {
	size, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("tail: seek end: %w", err)
	}
	if size == 0 {
		return nil, nil
	}

	var buf []byte
	pos := size
	linesNeeded := r.opts.NumLines + 1 // +1 to discard partial leading line

	for pos > 0 && bytes.Count(buf, []byte("\n")) < linesNeeded {
		chunk := r.opts.ChunkSize
		if pos < chunk {
			chunk = pos
		}
		pos -= chunk
		if _, err := rs.Seek(pos, io.SeekStart); err != nil {
			return nil, fmt.Errorf("tail: seek: %w", err)
		}
		tmp := make([]byte, chunk)
		n, err := io.ReadFull(rs, tmp)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("tail: read chunk: %w", err)
		}
		buf = append(tmp[:n], buf...)
	}

	scanner := bufio.NewScanner(bytes.NewReader(buf))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scan: %w", err)
	}

	// If we read from the middle of the file, drop the first partial line.
	if pos > 0 && len(lines) > 0 {
		lines = lines[1:]
	}
	if len(lines) > r.opts.NumLines {
		lines = lines[len(lines)-r.opts.NumLines:]
	}
	return lines, nil
}
