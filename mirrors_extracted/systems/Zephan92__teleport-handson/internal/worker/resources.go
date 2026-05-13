//go:build linux

package worker

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// getResourceLimits constructs the resource constraints for a new job.
// It combines hardcoded defaults for CPU/Memory with dynamic IO limits
// based on the block device backing the job's working directory.
func getResourceLimits(path string) (ResourceLimits, error) {
	// Note: These limits are currently hardcoded based on the challenge requirements.
	// In a production system, these would likely be passed in via the API request.
	ioLimit, err := getLinuxIOLimit(path)
	if err != nil {
		return ResourceLimits{}, fmt.Errorf("failed to calculate resource limits: %w", err)
	}

	return ResourceLimits{
		IOLimit:          ioLimit,
		MemoryLimitBytes: DefaultMemoryLimitBytes,
		CPULimitPercent:  DefaultCPULimitPercent,
	}, nil
}

// getLinuxIOLimit determines the block device major/minor numbers for the filesystem
// containing the given path and returns the default IO throughput limits for that device.
func getLinuxIOLimit(path string) (LinuxIOLimit, error) {
	major, minor, err := getDeviceMajorMinor(path)
	if err != nil {
		return LinuxIOLimit{}, fmt.Errorf("failed to resolve device for path %q: %w", path, err)
	}

	return LinuxIOLimit{
		ReadBPS:  DefaultDiskReadBPS,
		WriteBPS: DefaultDiskWriteBPS,
		Major:    major,
		Minor:    minor,
	}, nil
}

// getDeviceMajorMinor returns the major and minor device numbers for the
// filesystem containing the file at the specified path.
func getDeviceMajorMinor(path string) (uint32, uint32, error) {
	// Note: unix.Stat follows symlinks. This is intentional; if the working directory
	// is a symlink, we want the device ID of the actual storage location.
	var stat unix.Stat_t
	if err := unix.Stat(path, &stat); err != nil {
		return 0, 0, fmt.Errorf("failed to stat path %q: %w", path, err)
	}

	return unix.Major(uint64(stat.Dev)), unix.Minor(uint64(stat.Dev)), nil
}
