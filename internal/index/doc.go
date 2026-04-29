// Package index implements byte-offset indexing for log files.
//
// The index maps timestamps to byte offsets within a log file, allowing
// the scanner to seek directly to the approximate start of a time range
// rather than reading from the beginning of the file.
//
// Usage:
//
//	// Build an index by sampling every 100th line
//	builder := index.NewBuilder(100)
//	idx, err := builder.Build(file)
//
//	// Find the start offset for a time range
//	startOffset := idx.FindStart(startTime)
//	file.Seek(startOffset, io.SeekStart)
//
// The sample rate controls the trade-off between index size and seek
// precision. A lower sample rate produces a more precise index but uses
// more memory. For most use cases, sampling every 50-200 lines works well.
package index
