// Package sampler provides deterministic line sampling for logslice output.
//
// When processing very large log files it is sometimes useful to inspect a
// representative subset of matching lines rather than every line. The Sampler
// type supports two complementary strategies:
//
//   - Rate-based sampling: emit every Nth line (Every option). For example,
//     Every=10 emits roughly 10 % of matching lines while preserving their
//     relative order.
//
//   - Volume cap: stop emitting after MaxLines total lines have been produced,
//     regardless of the sampling rate. Zero disables the cap.
//
// Both strategies can be combined. The Sampler is not safe for concurrent use;
// callers that process lines from multiple goroutines should create one Sampler
// per goroutine or protect access with a mutex.
package sampler
