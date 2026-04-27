// Package scanner provides functionality for scanning large log files
// efficiently using binary search to locate time-based boundaries.
package scanner

import (
	"bufio"
	"io"
	"os"

	"github.com/yourorg/logslice/internal/logline"
)

// Scanner wraps a log file and provides methods for scanning lines
// within a given byte range, optionally filtered by a logline.Parser.
type Scanner struct {
	f      *os.File
	parser *logline.Parser
}

// New opens the file at path and returns a Scanner ready for use.
// The caller is responsible for calling Close when done.
func New(path string, parser *logline.Parser) (*Scanner, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Scanner{f: f, parser: parser}, nil
}

// Close releases the underlying file handle.
func (s *Scanner) Close() error {
	return s.f.Close()
}

// FileSize returns the total byte size of the underlying file.
func (s *Scanner) FileSize() (int64, error) {
	info, err := s.f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// LineAt seeks to offset and reads forward until a complete line is found,
// returning the parsed LogLine and the byte offset where the line starts.
// If offset is not at the beginning of a line, the next full line is returned.
func (s *Scanner) LineAt(offset int64) (logline.LogLine, int64, error) {
	if _, err := s.f.Seek(offset, io.SeekStart); err != nil {
		return logline.LogLine{}, 0, err
	}

	r := bufio.NewReader(s.f)

	// If we're not at the start of the file, skip the partial line we may
	// have landed in the middle of.
	if offset > 0 {
		if _, err := r.ReadString('\n'); err != nil && err != io.EOF {
			return logline.LogLine{}, 0, err
		}
		// Track how many bytes were consumed to align the offset.
		buffered := int64(r.Buffered())
		current, err := s.f.Seek(0, io.SeekCurrent)
		if err != nil {
			return logline.LogLine{}, 0, err
		}
		offset = current - buffered
	}

	for {
		lineStart := offset
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return logline.LogLine{}, 0, err
		}
		if line == "" {
			return logline.LogLine{}, 0, io.EOF
		}

		ll := s.parser.Parse(line, lineStart)
		offset += int64(len(line))

		if !ll.IsZero() {
			return ll, lineStart, nil
		}
		if err == io.EOF {
			return logline.LogLine{}, 0, io.EOF
		}
	}
}

// ScanRange reads all lines from startOffset to endOffset (inclusive byte
// range) and sends each raw line string to the out channel. The caller
// must close or drain the channel; ScanRange closes out when done.
func (s *Scanner) ScanRange(startOffset, endOffset int64, out chan<- string) error {
	defer close(out)

	if _, err := s.f.Seek(startOffset, io.SeekStart); err != nil {
		return err
	}

	limitedReader := io.LimitReader(s.f, endOffset-startOffset)
	scanner := bufio.NewScanner(limitedReader)

	for scanner.Scan() {
		out <- scanner.Text()
	}
	return scanner.Err()
}
