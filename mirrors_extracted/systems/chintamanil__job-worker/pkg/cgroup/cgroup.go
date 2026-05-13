//go:build linux

// Package cgroup provides cgroup v2 management for job process isolation.
//
// This package only supports Linux as cgroups are a Linux-specific feature.
package cgroup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	// ErrInvalidCPUQuota is returned when CPUQuota is negative.
	ErrInvalidCPUQuota = errors.New("cpu quota must be non-negative")

	// ErrInvalidMemoryLimit is returned when MemoryBytes is negative.
	ErrInvalidMemoryLimit = errors.New("memory limit must be non-negative")

	// ErrInvalidIOWeight is returned when IOWeight is out of range.
	ErrInvalidIOWeight = errors.New("io weight must be 0 or between 1 and 10000")

	// ErrNilCgroup is returned when a method is called on a nil Cgroup.
	ErrNilCgroup = errors.New("nil cgroup")
)

const (
	// cpuPeriod is the CPU bandwidth period in microseconds (100ms).
	cpuPeriod = 100000

	// Controller file names
	cpuMaxFile    = "cpu.max"
	memoryMaxFile = "memory.max"
	ioWeightFile  = "io.weight"
	cgroupKill    = "cgroup.kill"
	cgroupEvents  = "cgroup.events"

	// I/O weight limits
	minIOWeight = 1
	maxIOWeight = 10000
)

// Limits holds resource control settings for a cgroup.
type Limits struct {
	// CPUQuota is the fraction of CPU cores allowed (0.5 = 50%, 2.0 = 200%).
	// Zero means unlimited.
	CPUQuota float64

	// MemoryBytes is the hard memory limit in bytes.
	// Zero means unlimited. Exceeding triggers OOM killer.
	MemoryBytes int64

	// IOWeight is the relative I/O priority (1-10000, default 100).
	// Zero means use system default.
	IOWeight int32
}

// Validate checks that all limit values are within valid ranges.
func (l Limits) Validate() error {
	if l.CPUQuota < 0 {
		return ErrInvalidCPUQuota
	}
	if l.MemoryBytes < 0 {
		return ErrInvalidMemoryLimit
	}
	if l.IOWeight != 0 && (l.IOWeight < minIOWeight || l.IOWeight > maxIOWeight) {
		return ErrInvalidIOWeight
	}
	return nil
}

// Cgroup manages a single cgroup v2 for one job.
type Cgroup struct {
	id   string // Job ID (cgroup name)
	path string // Full path: <rootPath>/<id>
	fd   int    // Open directory FD, -1 if closed
}

// New creates a Cgroup instance without creating the actual cgroup.
// Call Create() to create the cgroup directory and set limits.
func New(rootPath, id string) *Cgroup {
	return &Cgroup{
		id:   id,
		path: filepath.Join(rootPath, id),
		fd:   -1,
	}
}

// Path returns the full cgroup directory path.
// Returns empty string if c is nil.
func (c *Cgroup) Path() string {
	if c == nil {
		return ""
	}
	return c.path
}

// Create creates the cgroup directory and writes resource limits.
// Limits with zero values are not written (use cgroup defaults).
// If any limit write fails, the cgroup directory is removed (rollback).
// Returns ErrNilCgroup if c is nil.
func (c *Cgroup) Create(limits Limits) error {
	if c == nil {
		return ErrNilCgroup
	}

	// Validate limits first
	if err := limits.Validate(); err != nil {
		return err
	}

	// Create the cgroup directory
	if err := os.MkdirAll(c.path, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory: %w", err)
	}

	// Track whether we need to clean up on failure
	needsCleanup := true
	defer func() {
		if needsCleanup {
			os.Remove(c.path)
		}
	}()

	// Write CPU limit if specified
	if limits.CPUQuota > 0 {
		quota := int(limits.CPUQuota * cpuPeriod)
		content := fmt.Sprintf("%d %d", quota, cpuPeriod)
		if err := c.writeController(cpuMaxFile, content); err != nil {
			return err
		}
	}

	// Write memory limit if specified
	if limits.MemoryBytes > 0 {
		content := strconv.FormatInt(limits.MemoryBytes, 10)
		if err := c.writeController(memoryMaxFile, content); err != nil {
			return err
		}
	}

	// Write I/O weight if specified
	if limits.IOWeight > 0 {
		content := fmt.Sprintf("default %d", limits.IOWeight)
		if err := c.writeController(ioWeightFile, content); err != nil {
			return err
		}
	}

	// Success - don't clean up
	needsCleanup = false
	return nil
}

// writeController writes a value to a cgroup controller file.
func (c *Cgroup) writeController(filename, content string) error {
	path := filepath.Join(c.path, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	return nil
}

// OpenFD opens the cgroup directory and returns the file descriptor.
// The FD is used with SysProcAttr.CgroupFD for atomic process assignment.
// Must call CloseFD() after cmd.Start() returns.
// Returns ErrNilCgroup if c is nil.
func (c *Cgroup) OpenFD() (int, error) {
	if c == nil {
		return -1, ErrNilCgroup
	}
	if c.fd != -1 {
		return c.fd, nil // Already open
	}

	// O_PATH: Get FD without requiring read/write access to contents
	// O_DIRECTORY: Ensure it's a directory
	// O_CLOEXEC: Don't leak FD to child processes after exec
	fd, err := unix.Open(c.path, unix.O_PATH|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, fmt.Errorf("failed to open cgroup directory: %w", err)
	}

	c.fd = fd
	return fd, nil
}

// CloseFD closes the cgroup directory FD. Safe to call multiple times.
// Returns nil if c is nil (no-op).
func (c *Cgroup) CloseFD() error {
	if c == nil || c.fd == -1 {
		return nil // Nil or already closed
	}

	if err := unix.Close(c.fd); err != nil {
		return fmt.Errorf("failed to close cgroup fd: %w", err)
	}

	c.fd = -1
	return nil
}

// Kill terminates all processes in the cgroup by writing to cgroup.kill.
// This sends SIGKILL to all processes atomically.
// Returns ErrNilCgroup if c is nil.
func (c *Cgroup) Kill() error {
	if c == nil {
		return ErrNilCgroup
	}
	if err := c.writeController(cgroupKill, "1"); err != nil {
		return fmt.Errorf("failed to kill cgroup: %w", err)
	}
	return nil
}

// Remove deletes the cgroup directory.
// Only succeeds if no processes remain in the cgroup.
// Returns ErrNilCgroup if c is nil.
func (c *Cgroup) Remove() error {
	if c == nil {
		return ErrNilCgroup
	}
	if err := os.Remove(c.path); err != nil {
		return fmt.Errorf("failed to remove cgroup: %w", err)
	}
	return nil
}

// IsEmpty returns true if the cgroup has no processes.
// Reads the "populated" field from cgroup.events.
// Returns ErrNilCgroup if c is nil.
func (c *Cgroup) IsEmpty() (bool, error) {
	if c == nil {
		return false, ErrNilCgroup
	}
	data, err := os.ReadFile(filepath.Join(c.path, cgroupEvents))
	if err != nil {
		return false, fmt.Errorf("failed to read cgroup.events: %w", err)
	}
	// cgroup.events format:
	// populated 0
	// frozen 0
	return strings.Contains(string(data), "populated 0"), nil
}

// Cleanup performs best-effort cleanup of all cgroup resources.
// It attempts to close the FD and remove the directory, collecting all errors.
// Returns nil if c is nil (no-op).
func (c *Cgroup) Cleanup() error {
	if c == nil {
		return nil
	}

	var errs []error

	if err := c.CloseFD(); err != nil {
		errs = append(errs, fmt.Errorf("close fd: %w", err))
	}

	if err := c.Remove(); err != nil {
		errs = append(errs, fmt.Errorf("remove: %w", err))
	}

	return errors.Join(errs...)
}
