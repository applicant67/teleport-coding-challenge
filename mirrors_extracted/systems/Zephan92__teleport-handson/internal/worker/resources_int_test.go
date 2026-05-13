//go:build linux

package worker

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestJob_MemoryLimitEnforcement asserts that a process exceeding its cgroup
// memory limit is terminated by the kernel.
//
// NOTE: This test requires privileges to create real cgroups and may need to
// be run with `sudo`. It also compiles a helper binary on the fly.
func TestJob_MemoryLimitEnforcement(t *testing.T) {
	t.Parallel()
	// This test interacts with the real cgroup filesystem, so we skip it
	// if not running with root privileges.
	if os.Geteuid() != 0 {
		t.Skip("Skipping memory limit test; requires root privileges")
	}

	// Create the helper source code in a temp dir
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "main.go")
	hogSource := `
package main
func main() {
	// Allocate 100 MB, which is more than the 10 MB limit.
	_ = make([]byte, 100*1024*1024)
	// Keep running so the OOM killer has time to act.
	select {}
}`
	if err := os.WriteFile(srcFile, []byte(hogSource), 0644); err != nil {
		t.Fatalf("Failed to write helper source: %v", err)
	}

	// Compile the memory hog helper program.
	tmpDir := t.TempDir()
	hogBin := filepath.Join(tmpDir, "memoryhog")
	buildCmd := exec.Command("go", "build", "-o", hogBin, srcFile)
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to compile memoryhog helper: %v", err)
	}

	// Set a very low memory limit.
	memLimitBytes := int64(10 * 1024 * 1024) // 10 MB

	cgroupRoot := "/sys/fs/cgroup/teleport-worker-test"
	// Ensure the test root exists and has controllers enabled
	if err := os.MkdirAll(cgroupRoot, 0755); err == nil {
		// Best effort to enable controllers. If this fails, the test might fail later.
		if err := os.WriteFile(filepath.Join(cgroupRoot, "cgroup.subtree_control"), []byte("+cpu +memory +io"), 0644); err != nil {
			t.Logf("Failed to enable controllers: %v", err)
		}
	}
	job, err := NewJob(ExecutionSpecs{
		Command:          hogBin,
		WorkingDirectory: tmpDir,
	}, cgroupRoot)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}
	job.Limit.MemoryLimitBytes = memLimitBytes

	if err := job.Start(t.Context()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for the job to be killed by the OOM killer.
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job to be OOM-killed")
		case <-ticker.C:
			s := job.Status()
			if s.Status == StatusFailed {
				// The kernel OOM killer typically results in exit code 137 (SIGKILL).
				if s.ExitCode == 0 {
					t.Errorf("Expected a non-zero exit code for OOM-killed process, got 0")
				}
				return // Success
			}
			if s.Status == StatusCompleted {
				t.Fatal("Job completed successfully but should have been OOM-killed")
			}
		}
	}
}
