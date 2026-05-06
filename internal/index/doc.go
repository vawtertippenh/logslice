// Package index provides time-based indexing for large log files.
//
// # Overview
//
// An Index maps sampled timestamps to byte offsets within a log file, enabling
// efficient seeking without scanning the entire file from the beginning.
//
// # Components
//
//   - Index / Builder: core data structure and construction from a log stream.
//   - Searcher: binary-search helpers to locate the start/end offset for a
//     given time range within an Index.
//   - Cache / TTLCache: in-memory caches for Index values, with optional LRU
//     eviction (Cache) or time-to-live expiry (TTLCache).
//   - LRUEviction: doubly-linked-list LRU policy used by Cache.
//   - Merger: combines multiple Index values into one, deduplicating entries.
//   - Stats / Collector: runtime metrics (hit rate, duration, empty checks)
//     collected concurrently and readable at any time.
//
// # Usage
//
// Build an index with NewBuilder, seek into it with NewSearcher, and cache the
// result with NewCache or NewTTLCache to avoid rebuilding on repeated queries
// against the same file.
package index
