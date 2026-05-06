// Package rotate detects and handles log file rotation for logslice.
//
// Log rotation can happen in two ways:
//
//  1. Replacement – the file is moved aside and a new file is created at the
//     same path (e.g. logrotate copytruncate or rename). This is detected by
//     an inode or device change.
//
//  2. Truncation – the file is opened and truncated in-place. This is detected
//     when the current file size is smaller than the last observed size.
//
// Detector provides a one-shot check. Watcher wraps a Detector in a polling
// loop and emits events on a channel so callers can react without blocking.
package rotate
