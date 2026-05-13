package worker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
)

// mount | grep cgroup
// cgroup2 on /sys/fs/cgroup type cgroup2 (rw,nosuid,nodev,noexec,relatime,seclabel,nsdelegate,memory_recursiveprot)
// With v2 heirarchy, this is the most common path though administrator can configure it differently.
const cgroupPath = "/sys/fs/cgroup"

type ResourceLimits struct {
	CPULimit     int
	MemoryLimit  int
	DiskIOWeight int
}

type ResourceManager struct {
	limits map[string]ResourceLimits
	mu     sync.RWMutex
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		limits: make(map[string]ResourceLimits),
	}
}

func (rm *ResourceManager) SetLimits(jobID string, cpuLimit int, memoryLimit int, diskIOWeight int) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.limits[jobID] = ResourceLimits{
		CPULimit:     cpuLimit,
		MemoryLimit:  memoryLimit,
		DiskIOWeight: diskIOWeight,
	}
	return nil
}

func (rm *ResourceManager) GetLimits(jobID string) (int, int, int, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	limits, exists := rm.limits[jobID]
	if !exists {
		return 0, 0, 0, fmt.Errorf("job ID %s not found", jobID)
	}
	return limits.CPULimit, limits.MemoryLimit, limits.DiskIOWeight, nil
}

func (rm *ResourceManager) CreateCgroup(jobID string) error {
	cgroupName := fmt.Sprintf("job_%s", jobID)
	cgroupPath := filepath.Join(cgroupPath, cgroupName)
	err := os.Mkdir(cgroupPath, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("error creating cgroup: %w", err)
		}
		// Cgroup already exists, continue
	}

	return nil
}

// Additionally we might need a cronjob to clean up cgroups
// Example: rmdir job_* inside /sys/fs/cgroup folder
func (rm *ResourceManager) CleanupCgroup(jobID string) error {
	cgroupName := fmt.Sprintf("job_%s", jobID)
	cgroupPath := filepath.Join(cgroupPath, cgroupName)
	err := os.RemoveAll(cgroupPath)
	if err != nil {
		return fmt.Errorf("error removing cgroup: %w", err)
	}
	return nil
}

func (rm *ResourceManager) StartProcessInCgroup(jobID string, cmd *exec.Cmd) error {
	cgroupName := fmt.Sprintf("job_%s", jobID)
	cgroupDir := filepath.Join(cgroupPath, cgroupName)

	cgroupFD, err := os.Open(cgroupDir)
	if err != nil {
		return fmt.Errorf("error opening cgroup directory: %w", err)
	}
	defer cgroupFD.Close()

	// Set resource limits for the job
	cpuLimit, memoryLimit, diskIOWeight, err := rm.GetLimits(jobID)
	if err != nil {
		return fmt.Errorf("error getting resource limits: %w", err)
	}

	// set memory limit
	err = os.WriteFile(filepath.Join(cgroupDir, "memory.max"), []byte(strconv.Itoa(memoryLimit)), 0644)
	if err != nil {
		return fmt.Errorf("error setting memory limit: %w", err)
	}

	// Set CPU limit
	err = os.WriteFile(filepath.Join(cgroupDir, "cpu.max"), []byte(fmt.Sprintf("%d 100000", cpuLimit)), 0644)
	if err != nil {
		return fmt.Errorf("error setting CPU limit: %w", err)
	}

	// Set disk I/O limit
	err = os.WriteFile(filepath.Join(cgroupDir, "io.weight"), []byte(strconv.Itoa(diskIOWeight)), 0644)
	if err != nil {
		return fmt.Errorf("failed to set disk I/O limit: %w", err)
	}

	// Set the SysProcAttr to specify cgroup settings for the new process before starting it
	// added in go 1.22 plus. https://pkg.go.dev/syscall#SysProcAttr
	// some discussions that were helpful in grokking what is going on: https://github.com/golang/go/issues/51246
	// Alternative method seems to be using the CLONE_INTO_CGROUP flag which does unix.Clone
	// Clone can 1> Create a new process 2> control behavior of child process.
	// Under the hood, UseCgroupFD seems to be doing the same as per here: https://cs.opensource.google/go/go/+/master:src/syscall/exec_linux.go;l=312
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Enable using the cgroup file descriptor
		UseCgroupFD: true,
		// Set the cgroup file descriptor to the opened cgroup directory
		CgroupFD: int(cgroupFD.Fd()),
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	return nil
}
