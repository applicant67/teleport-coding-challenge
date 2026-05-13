//go:build linux

package worker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestJob_ProcessTreeCleanup asserts that stopping a job cleans up the entire
// process tree, including child processes spawned by the main command.
func TestJob_ProcessTreeCleanup(t *testing.T) {
	t.Parallel()
	if os.Geteuid() != 0 {
		t.Skip("Skipping process tree cleanup test; requires root privileges for cgroups")
	}

	tmpDir := t.TempDir()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// Create a script that spawns background processes
	scriptPath := filepath.Join(tmpDir, "tree.sh")
	scriptContent := `#!/bin/sh
sleep 10 &
sleep 10 &
wait
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to write script: %v", err)
	}

	cgroupRoot := "/sys/fs/cgroup/teleport-worker-test"
	// Ensure the test root exists and has controllers enabled
	if err := os.MkdirAll(cgroupRoot, 0755); err == nil {
		// Best effort to enable controllers. If this fails, the test might fail later.
		if err := os.WriteFile(filepath.Join(cgroupRoot, "cgroup.subtree_control"), []byte("+cpu +memory +io"), 0644); err != nil {
			t.Logf("Failed to enable controllers: %v", err)
		}
	}
	job, err := NewJob(ExecutionSpecs{
		Command:          scriptPath,
		WorkingDirectory: tmpDir,
	}, cgroupRoot)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for the process tree to populate
	cgroupProcs := filepath.Join(cgroupRoot, string(job.ID), "cgroup.procs")
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	seenMultiple := false
	for !seenMultiple {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for process tree to populate")
		case <-ticker.C:
			content, err := os.ReadFile(cgroupProcs)
			if err == nil {
				// PIDs are newline separated
				pids := strings.Fields(string(content))
				if len(pids) >= 2 {
					seenMultiple = true
					break
				}
			}
		}
	}

	// Stop the job
	if err := job.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Verify cleanup: The cgroup directory should be gone.
	// If child processes were not killed, the cgroup would remain (EBUSY on deletion).
	cgroupPath := filepath.Join(cgroupRoot, string(job.ID))
	for i := 0; i < 50; i++ {
		if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
			return // Success
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Errorf("Cgroup directory %s still exists after Stop()", cgroupPath)
}
