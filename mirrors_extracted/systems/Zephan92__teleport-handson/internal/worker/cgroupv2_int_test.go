//go:build linux

package worker

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// TestJob_CgroupIntegration asserts that a Job correctly uses its Cgroup
// controller to create the cgroup directory and apply the specified resource
// limits to the filesystem.
func TestJob_CgroupIntegration(t *testing.T) {
	t.Parallel()
	// Override CgroupRoot to a temporary directory for testing
	tmpDir := t.TempDir()

	job, err := NewJob(ExecutionSpecs{
		Command:          "echo",
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	var stat syscall.Stat_t
	if err := syscall.Stat(".", &stat); err != nil {
		t.Fatalf("failed to stat .: %v", err)
	}

	job.Limit = ResourceLimits{
		MemoryLimitBytes: 1024 * 1024, // 1MB
		CPULimitPercent:  0.5,
		IOLimit: LinuxIOLimit{
			Major:    unix.Major(stat.Dev),
			Minor:    unix.Minor(stat.Dev),
			ReadBPS:  1024 * 1024,
			WriteBPS: 1024 * 1024,
		},
	}

	if err := job.Cgroup.Create(job.Limit); err != nil {
		t.Fatalf("Cgroup.Create failed: %v", err)
	}

	// Verify directory exists
	cgroupPath := filepath.Join(tmpDir, string(job.ID))
	if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
		t.Errorf("cgroup directory was not created at %s", cgroupPath)
	}

	// Verify memory limit
	assertFileContent(t, filepath.Join(cgroupPath, cgroupMemoryMaxFile), "1048576")

	// Verify CPU limit
	assertFileContent(t, filepath.Join(cgroupPath, cgroupCPUMaxFile), "50000 100000")
}

// TestCgroup_Attach_InvalidPID asserts that Attach fails if the PID does not exist
// when running against a real cgroup hierarchy.
func TestCgroup_Attach_InvalidPID(t *testing.T) {
	t.Parallel()
	if os.Geteuid() != 0 {
		t.Skip("Skipping cgroup attach test; requires root privileges")
	}

	cgroupRoot := "/sys/fs/cgroup/teleport-worker-test-attach"
	if err := os.MkdirAll(cgroupRoot, 0755); err != nil {
		t.Fatalf("Failed to create cgroup root: %v", err)
	}
	// Enable controllers (best effort)
	if err := os.WriteFile(filepath.Join(cgroupRoot, "cgroup.subtree_control"), []byte("+cpu +memory +io"), 0644); err != nil {
		t.Logf("Failed to enable controllers: %v", err)
	}
	defer os.RemoveAll(cgroupRoot)

	cg := NewCgroup("test-invalid-pid", cgroupRoot)

	if err := cg.Create(ResourceLimits{
		CPULimitPercent:  1.0,
		MemoryLimitBytes: 1024 * 1024,
	}); err != nil {
		t.Fatalf("Failed to create cgroup: %v", err)
	}
	defer cg.Delete()

	if err := cg.Attach(99999999); err == nil {
		t.Error("Expected error attaching invalid PID to real cgroup")
	}
}
