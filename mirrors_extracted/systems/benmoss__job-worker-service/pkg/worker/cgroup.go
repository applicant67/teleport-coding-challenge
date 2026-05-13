package worker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const defaultCgroupRoot = "/sys/fs/cgroup/jobworker"

// cgroup manages a cgroups v2 cgroup for a job.
type cgroup struct {
	path string
}

// createCgroupInput contains parameters for creating a new cgroup.
type createCgroupInput struct {
	cgroupRoot     string
	jobID          uint64
	resourceLimits *ResourceLimits
}

// createCgroup creates a new cgroup for a job and applies resource limits.
func createCgroup(input createCgroupInput) (*cgroup, error) {
	root := input.cgroupRoot
	if root == "" {
		root = defaultCgroupRoot
	}

	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create cgroup root %s: %w", root, err)
	}

	// Enable controllers in parent cgroup
	if err := enableControllers(root); err != nil {
		return nil, fmt.Errorf("enable controllers in %s: %w", root, err)
	}

	cgPath := filepath.Join(root, fmt.Sprintf("job-%d", input.jobID))

	if err := os.MkdirAll(cgPath, 0o755); err != nil {
		return nil, fmt.Errorf("create cgroup dir %s: %w", cgPath, err)
	}

	cg := &cgroup{path: cgPath}

	if err := cg.applyLimits(input.resourceLimits); err != nil {
		return nil, fmt.Errorf("apply resource limits: %w", err)
	}

	return cg, nil
}

// cgroupFD is a cgroup file descriptor that must be closed after use.
type cgroupFD struct {
	fd int
}

// Close closes the file descriptor.
func (c *cgroupFD) Close() error {
	return syscall.Close(c.fd)
}

// FD returns the raw file descriptor value.
func (c *cgroupFD) FD() int {
	return c.fd
}

// fd returns a file descriptor for the cgroup directory for use with
// CLONE_INTO_CGROUP via exec.Cmd.SysProcAttr.CgroupFD.
// The caller must call Close() on the returned cgroupFD when done.
func (c *cgroup) fd() (*cgroupFD, error) {
	fd, err := syscall.Open(c.path, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, fmt.Errorf("open cgroup dir %s: %w", c.path, err)
	}
	return &cgroupFD{fd}, nil
}

// kill terminates all processes in the cgroup by writing to cgroup.kill.
func (c *cgroup) kill() error {
	killPath := filepath.Join(c.path, "cgroup.kill")
	if err := os.WriteFile(killPath, []byte("1"), 0o644); err != nil {
		return fmt.Errorf("write cgroup.kill: %w", err)
	}
	return nil
}

// isPopulated returns true if the cgroup contains any processes.
// It parses the cgroup.events file and checks for "populated 1".
func (c *cgroup) isPopulated() (bool, error) {
	eventsPath := filepath.Join(c.path, "cgroup.events")
	content, err := os.ReadFile(eventsPath)
	if err != nil {
		return false, fmt.Errorf("read cgroup.events: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "populated ") {
			fields := strings.Fields(line)
			if len(fields) == 2 && fields[1] == "1" {
				return true, nil
			}
			return false, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("scan cgroup.events: %w", err)
	}

	return false, fmt.Errorf("populated field not found in cgroup.events")
}

// waitUntilEmpty polls cgroup.events until populated=0.
// This waits indefinitely until all processes in the cgroup have exited,
// similar to how cmd.Wait() waits for a process.
func (c *cgroup) waitUntilEmpty() error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		populated, err := c.isPopulated()
		if err != nil {
			return fmt.Errorf("check if cgroup populated: %w", err)
		}

		if !populated {
			return nil
		}

		<-ticker.C
	}
}

// remove deletes the cgroup directory.
func (c *cgroup) remove() error {
	if err := os.Remove(c.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove cgroup dir %s: %w", c.path, err)
	}
	return nil
}

// applyLimits writes resource limits to the cgroup control files.
func (c *cgroup) applyLimits(limits *ResourceLimits) error {
	if limits == nil {
		return nil
	}

	// Apply CPU limit
	if limits.CPUMax > 0 {
		if err := c.applyCPULimit(limits.CPUMax); err != nil {
			return fmt.Errorf("apply cpu limit: %w", err)
		}
	}

	// Apply memory limit
	if limits.MemoryMaxBytes > 0 {
		if err := c.applyMemoryLimit(limits.MemoryMaxBytes); err != nil {
			return fmt.Errorf("apply memory limit: %w", err)
		}
	}

	// Apply I/O limits
	if limits.IOMaxReadBPS > 0 || limits.IOMaxWriteBPS > 0 {
		if err := c.applyIOLimits(limits.IOMaxReadBPS, limits.IOMaxWriteBPS); err != nil {
			return fmt.Errorf("apply io limits: %w", err)
		}
	}

	return nil
}

// applyCPULimit sets the CPU bandwidth limit.
// cpuMax is a fraction of CPU (0.5 = half a core, 2.0 = two cores).
// This is converted to cgroup cpu.max format: "quota period" in microseconds.
func (c *cgroup) applyCPULimit(cpuMax float64) error {
	const period = 100000 // 100ms in microseconds
	quota := int64(cpuMax * float64(period))

	cpuMaxPath := filepath.Join(c.path, "cpu.max")
	content := fmt.Sprintf("%d %d", quota, period)

	if err := os.WriteFile(cpuMaxPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", cpuMaxPath, err)
	}

	return nil
}

// applyMemoryLimit sets the memory limit.
func (c *cgroup) applyMemoryLimit(memoryMaxBytes int64) error {
	memoryMaxPath := filepath.Join(c.path, "memory.max")
	content := fmt.Sprintf("%d", memoryMaxBytes)

	if err := os.WriteFile(memoryMaxPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", memoryMaxPath, err)
	}

	return nil
}

// applyIOLimits sets I/O bandwidth limits for all block devices.
func (c *cgroup) applyIOLimits(readBPS, writeBPS int64) error {
	devices, err := getBlockDevices()
	if err != nil {
		return fmt.Errorf("get block devices: %w", err)
	}

	if len(devices) == 0 {
		return nil
	}

	ioMaxPath := filepath.Join(c.path, "io.max")
	var lines []string

	for _, dev := range devices {
		var parts []string
		parts = append(parts, fmt.Sprintf("%d:%d", dev.major, dev.minor))

		if readBPS > 0 {
			parts = append(parts, fmt.Sprintf("rbps=%d", readBPS))
		}
		if writeBPS > 0 {
			parts = append(parts, fmt.Sprintf("wbps=%d", writeBPS))
		}

		if len(parts) > 1 {
			lines = append(lines, strings.Join(parts, " "))
		}
	}

	if len(lines) == 0 {
		return nil
	}

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(ioMaxPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", ioMaxPath, err)
	}

	return nil
}

// blockDevice represents a block device with major and minor numbers.
type blockDevice struct {
	major uint64
	minor uint64
}

// getBlockDevices returns a list of block devices from /proc/partitions.
func getBlockDevices() ([]blockDevice, error) {
	file, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, fmt.Errorf("open /proc/partitions: %w", err)
	}
	defer file.Close()

	var devices []blockDevice
	scanner := bufio.NewScanner(file)

	// Skip header lines
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "major") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		major, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		minor, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		devices = append(devices, blockDevice{major: major, minor: minor})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan /proc/partitions: %w", err)
	}

	return devices, nil
}

// enableControllers enables cgroup controllers needed for resource limits.
func enableControllers(cgroupPath string) error {
	subtreeControlPath := filepath.Join(cgroupPath, "cgroup.subtree_control")

	// Read current controllers
	current, err := os.ReadFile(subtreeControlPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read %s: %w", subtreeControlPath, err)
	}

	// Parse currently enabled controllers
	currentControllers := make(map[string]bool)
	for _, ctrl := range strings.Fields(string(current)) {
		currentControllers[ctrl] = true
	}

	needed := []string{"cpu", "memory", "io"}
	var toEnable []string

	for _, controller := range needed {
		if !currentControllers[controller] {
			toEnable = append(toEnable, "+"+controller)
		}
	}

	if len(toEnable) == 0 {
		return nil
	}

	content := strings.Join(toEnable, " ")
	if err := os.WriteFile(subtreeControlPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", subtreeControlPath, err)
	}

	return nil
}

// cleanupCgroups removes all existing job cgroups under the root.
func cleanupCgroups(cgroupRoot string) error {
	if cgroupRoot == "" {
		cgroupRoot = defaultCgroupRoot
	}

	entries, err := os.ReadDir(cgroupRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read cgroup root %s: %w", cgroupRoot, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "job-") {
			continue
		}

		cgPath := filepath.Join(cgroupRoot, entry.Name())
		cg := &cgroup{path: cgPath}

		if err := cg.remove(); err != nil {
			return fmt.Errorf("remove cgroup %s: %w", cgPath, err)
		}
	}

	return nil
}
