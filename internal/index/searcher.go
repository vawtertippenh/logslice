package index

import (
	"io"
	"time"
)

// Searcher uses an Index to efficiently seek to a start position
// in a ReadSeeker before scanning begins.
type Searcher struct {
	idx *Index
}

// NewSearcher returns a Searcher backed by the given Index.
func NewSearcher(idx *Index) *Searcher {
	return &Searcher{idx: idx}
}

// SeekToStart seeks rs to the byte offset of the last index entry whose
// timestamp is <= start, so that scanning begins as close to the desired
// time range as possible without skipping any matching lines.
// If no suitable entry is found, the reader is seeked to the beginning.
// Returns the offset that was seeked to.
func (s *Searcher) SeekToStart(rs io.ReadSeeker, start time.Time) (int64, error) {
	offset := s.idx.FindStart(start)
	_, err := rs.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

// SeekToEnd seeks rs to the byte offset of the first index entry whose
// timestamp is > end, so that scanning can stop early.
// Returns -1 if no such entry exists (scan until EOF).
// The reader position is NOT changed; callers use the returned offset
// as a read limit.
func (s *Searcher) SeekToEnd(rs io.ReadSeeker, end time.Time) int64 {
	return s.idx.FindEnd(end)
}

// SeekToStartAndEnd is a convenience wrapper that returns both the start
// offset (seeked to in rs) and the end offset limit.
func (s *Searcher) SeekToStartAndEnd(rs io.ReadSeeker, start, end time.Time) (startOffset, endOffset int64, err error) {
	startOffset, err = s.SeekToStart(rs, start)
	if err != nil {
		return 0, -1, err
	}
	endOffset = s.SeekToEnd(rs, end)
	return startOffset, endOffset, nil
}
