//go:build linux

package worker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

const (
	// CgroupMountPoint is the standard mount point for cgroup v2.
	CgroupMountPoint = "/sys/fs/cgroup"
	// CgroupWorkerDir is the name of the subdirectory for worker cgroups.
	CgroupWorkerDir = "worker"
	// dirMode is the permission bits for directories (rwxr-xr-x).
	dirMode = 0755
	// fileMode is the permission bits for files (rw-r--r--).
	fileMode = 0644
	// cpuPeriodUS is the standard CFS period in microseconds (100ms).
	cpuPeriodUS = 100000

	// Cgroup filenames
	cgroupMemoryMaxFile      = "memory.max"
	cgroupCPUMaxFile         = "cpu.max"
	cgroupIOMaxFile          = "io.max"
	cgroupProcsFile          = "cgroup.procs"
	cgroupKillFile           = "cgroup.kill"
	cgroupSubtreeControlFile = "cgroup.subtree_control"
	cgroupControllersFile    = "cgroup.controllers"
)

// Cgroup manages the Linux Control Group for a specific job.
type Cgroup struct {
	// JobID is the unique identifier for the job this cgroup belongs to.
	JobID JobID
	// Root is the base directory for cgroups (e.g., /sys/fs/cgroup/worker).
	Root string
}

// NewCgroup creates a new Cgroup controller.
// If root is empty, it defaults to the system standard (/sys/fs/cgroup/worker).
func NewCgroup(id JobID, root string) *Cgroup {
	// If no root is provided, use the default system path.
	if root == "" {
		root = filepath.Join(CgroupMountPoint, CgroupWorkerDir)
	}
	return &Cgroup{
		JobID: id,
		Root:  root,
	}
}

// Create initializes the Cgroup v2 directory hierarchy and applies
// memory, CPU, and IO controller limits before process attachment.
//
// Prerequisite: The parent directory (CgroupRoot) must have the target controllers
// (cpu, memory, io) enabled in its cgroup.subtree_control file.
func (c *Cgroup) Create(limit ResourceLimits) error {
	cgroupPath := c.path()

	if err := os.MkdirAll(cgroupPath, dirMode); err != nil {
		return fmt.Errorf("failed to create cgroup dir: %w", err)
	}

	cpuMax := formatCPULimit(limit.CPULimitPercent)
	memMax := formatMemoryLimit(limit.MemoryLimitBytes)
	ioMax := limit.IOLimit.String()

	if err := writeCgroupFile(cgroupPath, cgroupMemoryMaxFile, memMax); err != nil {
		c.Delete()
		return fmt.Errorf("failed to set memory limit: %w", err)
	}
	if err := writeCgroupFile(cgroupPath, cgroupCPUMaxFile, cpuMax); err != nil {
		c.Delete()
		return fmt.Errorf("failed to set cpu limit: %w", err)
	}
	// Only apply IO limits if a valid block device is specified.
	// In production, Major and Minor are expected to be present.
	// We allow 0 values here to simplify testing where device lookup is brittle.
	if limit.IOLimit.Major != 0 || limit.IOLimit.Minor != 0 {
		if err := writeCgroupFile(cgroupPath, cgroupIOMaxFile, ioMax); err != nil {
			c.Delete()
			return fmt.Errorf("failed to set io limit: %w", err)
		}
	}
	return nil
}

// Attach moves the target PID into the job's dedicated cgroup.procs.
func (c *Cgroup) Attach(pid ProcessID) error {
	if err := writeCgroupFile(c.path(), cgroupProcsFile, pid.String()); err != nil {
		return fmt.Errorf("failed to attach PID to cgroup: %w", err)
	}
	return nil
}

// Delete cleans up the cgroup directory.
// Note: This operation may fail if the cgroup is not empty (i.e., if processes
// are still running or are zombies). Ensure Cgroup.Kill() and exec.Cmd.Wait() have completed first.
// We use unix.Rmdir directly because it is the only correct way to remove a cgroup.
// See: https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html#organizing-processes
func (c *Cgroup) Delete() error {
	// Retry deletion a few times to handle race conditions where the kernel
	// hasn't fully reaped the processes yet.
	var err error
	for i := 0; i < 5; i++ {
		err = unix.Rmdir(c.path())
		if err == nil || os.IsNotExist(err) {
			return nil
		}
		// Only retry if the error is EBUSY, otherwise fail fast.
		if !errors.Is(err, unix.EBUSY) {
			if errors.Is(err, unix.ENOTEMPTY) && c.handleTestDirCleanup() {
				return nil
			}
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	return err
}

// Kill triggers the cgroup.kill mechanism to terminate all processes in the group.
// This requires Linux kernel 5.14+.
func (c *Cgroup) Kill() error {
	// Writing "1" to cgroup.kill kills every process in the cgroup immediately.
	return writeCgroupFile(c.path(), cgroupKillFile, "1")
}

// path returns the absolute path to the cgroup directory.
func (c *Cgroup) path() string {
	return filepath.Join(c.Root, string(c.JobID))
}

func formatCPULimit(percent float64) string {
	// Format: $QUOTA $PERIOD (e.g., 50% = 50000 100000)
	quota := int64(percent * float64(cpuPeriodUS))
	return fmt.Sprintf("%d %d", quota, cpuPeriodUS)
}

func formatMemoryLimit(bytes int64) string {
	return strconv.FormatInt(bytes, 10)
}

// writeCgroupFile writes the content to the specified file within the directory.
func writeCgroupFile(dir, file, content string) error {
	return os.WriteFile(filepath.Join(dir, file), []byte(content), fileMode)
}

// handleTestDirCleanup handles a specific edge case for integration tests.
//
// When running tests with t.TempDir(), the "cgroup" is just a regular directory on ext4/tmpfs.
// unix.Rmdir fails with ENOTEMPTY on regular directories if they contain files (like memory.max),
// whereas on a real cgroup v2 filesystem, Rmdir correctly removes the cgroup anchor even if
// controller interface files exist.
//
// This function checks the filesystem magic number. If it is NOT a cgroup filesystem,
// it assumes we are in a test environment and safely performs a recursive delete.
func (c *Cgroup) handleTestDirCleanup() bool {
	var stat unix.Statfs_t
	if err := unix.Statfs(c.path(), &stat); err != nil {
		return false
	}

	// Safety Check: If this IS a real cgroup filesystem, we must rely on Rmdir.
	// Recursive deletion on cgroupfs is undefined/dangerous.
	if stat.Type == unix.CGROUP2_SUPER_MAGIC {
		return false
	}

	// Safe to use RemoveAll for test directories
	return os.RemoveAll(c.path()) == nil
}
