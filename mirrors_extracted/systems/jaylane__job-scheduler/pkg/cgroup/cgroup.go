package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
)

const cgroupRoot = "/sys/fs/cgroup/"

// CreateCgroup creates a cgroup with the given path and sets the memory limit, CPU shares, CPU period, CPU quota, and disk I/O weight.
func CreateCgroup(groupPath string) error {
	// Create the cgroup directory if it doesn't exist
	if err := os.MkdirAll(fmt.Sprintf("%s%s", cgroupRoot, groupPath), 0755); err != nil {
		return fmt.Errorf("failed to create cgroup: %v", err)
	}

	// Set the memory limit (in bytes) for the cgroup
	if err := os.WriteFile(filepath.Join(groupPath, "memory.limit_in_bytes"), []byte("104857600"), 0644); err != nil {
		return fmt.Errorf("failed to set memory limit: %v", err)
	}

	// Set the CPU shares for the cgroup (number of CPU shares)
	if err := os.WriteFile(filepath.Join(groupPath, "cpu.shares"), []byte("2048"), 0644); err != nil {
		return fmt.Errorf("failed to set cpu shares: %v", err)
	}

	// Set the CPU period for the cgroup (in microseconds)
	if err := os.WriteFile(filepath.Join(groupPath, "cpu.period"), []byte("100000"), 0644); err != nil {
		return fmt.Errorf("failed to set cpu period: %v", err)
	}

	// Set the CPU quota for the cgroup (in microseconds)
	if err := os.WriteFile(filepath.Join(groupPath, "cpu.quota"), []byte("50000"), 0644); err != nil {
		return fmt.Errorf("failed to set cpu quota: %v", err)
	}

	// Set the disk I/O weight for the cgroup (1-10000, default is 1000)
	if err := os.WriteFile(filepath.Join(groupPath, "blkio.weight"), []byte("500"), 0644); err != nil {
		return fmt.Errorf("failed to set blkio weight: %v", err)
	}

	return nil
}
