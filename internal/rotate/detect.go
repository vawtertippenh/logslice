// Package rotate provides utilities for detecting and handling log file rotation.
// It monitors inode changes, truncation, and file replacement to support
// continuous tailing across log rotations.
package rotate

import (
	"os"
	"syscall"
)

// State captures the identity of a file at a point in time.
type State struct {
	Inode  uint64
	Device uint64
	Size   int64
}

// Detector checks whether a file has been rotated since it was last observed.
type Detector struct {
	path  string
	last  State
}

// NewDetector creates a Detector for the given file path and records its
// current state as the baseline.
func NewDetector(path string) (*Detector, error) {
	d := &Detector{path: path}
	st, err := d.stat()
	if err != nil {
		return nil, err
	}
	d.last = st
	return d, nil
}

// Rotated returns true if the file has been rotated (replaced or truncated)
// since the detector was created or last reset.
func (d *Detector) Rotated() (bool, error) {
	cur, err := d.stat()
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	if cur.Inode != d.last.Inode || cur.Device != d.last.Device {
		return true, nil
	}
	if cur.Size < d.last.Size {
		return true, nil
	}
	return false, nil
}

// Reset updates the baseline state to the current file state.
func (d *Detector) Reset() error {
	st, err := d.stat()
	if err != nil {
		return err
	}
	d.last = st
	return nil
}

func (d *Detector) stat() (State, error) {
	info, err := os.Stat(d.path)
	if err != nil {
		return State{}, err
	}
	sys, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return State{Size: info.Size()}, nil
	}
	return State{
		Inode:  sys.Ino,
		Device: uint64(sys.Dev),
		Size:   info.Size(),
	}, nil
}
