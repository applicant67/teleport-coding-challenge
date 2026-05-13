//go:build linux

package worker

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCgroup_Create asserts that the Create method correctly creates a cgroup
// directory and writes the formatted resource limits to the appropriate files.
func TestCgroup_Create(t *testing.T) {
	// Use a temporary directory as the cgroup root
	tmpDir := t.TempDir()

	jobID := JobID("test-job-123")
	cg := NewCgroup(jobID, tmpDir)

	limits := ResourceLimits{
		MemoryLimitBytes: 1024 * 1024, // 1MB
		CPULimitPercent:  0.5,         // 50%
		IOLimit: LinuxIOLimit{
			Major:    8,
			Minor:    0,
			ReadBPS:  1024,
			WriteBPS: 2048,
		},
	}

	// Test Create
	if err := cg.Create(limits); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify directory creation
	cgroupPath := filepath.Join(tmpDir, string(jobID))
	if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
		t.Errorf("Cgroup directory not created at %s", cgroupPath)
	}

	// Verify file contents
	assertFileContent(t, filepath.Join(cgroupPath, "memory.max"), "1048576")
	assertFileContent(t, filepath.Join(cgroupPath, "cpu.max"), "50000 100000") // 50% of 100000
	assertFileContent(t, filepath.Join(cgroupPath, "io.max"), "8:0 rbps=1024 wbps=2048")
}

// TestCgroup_Attach asserts that the Attach method correctly writes a given
// PID to the cgroup.procs file.
func TestCgroup_Attach(t *testing.T) {
	tmpDir := t.TempDir()

	jobID := JobID("test-job-attach")
	cg := NewCgroup(jobID, tmpDir)

	// Manually create the directory (simulating Create called previously)
	cgroupPath := filepath.Join(tmpDir, string(jobID))
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		t.Fatalf("Failed to create mock cgroup dir: %v", err)
	}

	pid := ProcessID(12345)
	if err := cg.Attach(pid); err != nil {
		t.Fatalf("Attach() failed: %v", err)
	}

	assertFileContent(t, filepath.Join(cgroupPath, "cgroup.procs"), "12345")
}

// TestCgroup_Kill asserts that the Kill method writes "1" to the cgroup.kill
// file to terminate all processes in the group.
func TestCgroup_Kill(t *testing.T) {
	tmpDir := t.TempDir()

	jobID := JobID("test-job-kill")
	cg := NewCgroup(jobID, tmpDir)
	cgroupPath := filepath.Join(tmpDir, string(jobID))
	os.MkdirAll(cgroupPath, 0755)

	if err := cg.Kill(); err != nil {
		t.Fatalf("Kill() failed: %v", err)
	}

	assertFileContent(t, filepath.Join(cgroupPath, "cgroup.kill"), "1")
}

// TestCgroup_Delete asserts that the Delete method successfully removes the
// cgroup directory from the filesystem.
func TestCgroup_Delete(t *testing.T) {
	tmpDir := t.TempDir()

	jobID := JobID("test-job-delete")
	cg := NewCgroup(jobID, tmpDir)
	cgroupPath := filepath.Join(tmpDir, string(jobID))
	os.MkdirAll(cgroupPath, 0755)

	if err := cg.Delete(); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	if _, err := os.Stat(cgroupPath); !os.IsNotExist(err) {
		t.Errorf("Cgroup directory was not deleted")
	}
}

func assertFileContent(t *testing.T, path, expected string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", path, err)
	}
	if string(content) != expected {
		t.Errorf("File %s content mismatch.\nGot:  %q\nWant: %q", path, string(content), expected)
	}
}

// TestCgroup_Errors covers failure scenarios for Cgroup operations.
func TestCgroup_Errors(t *testing.T) {
	// Create_InvalidRoot asserts that Create fails if the root directory cannot be created.
	t.Run("Create_InvalidRoot", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Make root a file to cause MkdirAll to fail
		dummyFile := filepath.Join(tmpDir, "file")
		_ = os.WriteFile(dummyFile, []byte(""), 0644)
		cg := NewCgroup("test-err", dummyFile)
		if err := cg.Create(ResourceLimits{}); err == nil {
			t.Error("Expected error when root is a file")
		}
	})

	// Attach_NonExistent asserts that Attach fails if the cgroup directory does not exist.
	t.Run("Attach_NonExistent", func(t *testing.T) {
		tmpDir := t.TempDir()

		cg := NewCgroup("test-attach-err", tmpDir)
		// Do not Create(), so directory doesn't exist

		if err := cg.Attach(123); err == nil {
			t.Error("Expected error attaching to non-existent cgroup")
		}
	})

	// Kill_NonExistent asserts that Kill fails if the cgroup directory does not exist.
	t.Run("Kill_NonExistent", func(t *testing.T) {
		tmpDir := t.TempDir()

		cg := NewCgroup("test-kill-err", tmpDir)
		// Do not Create(), so directory doesn't exist

		// Kill tries to write to cgroup.kill
		if err := cg.Kill(); err == nil {
			t.Error("Expected error killing non-existent cgroup")
		}
	})

	// Create_PartialWriteFailure asserts that Create cleans up the cgroup directory
	// if one of the file writes fails midway through.
	t.Run("Create_PartialWriteFailure", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("Skipping test that relies on permission errors (running as root)")
		}
		tmpDir := t.TempDir()

		cg := NewCgroup("test-partial-fail", tmpDir)
		cgroupPath := cg.path()

		// Pre-create the directory but make it read-only to cause a write failure.
		if err := os.MkdirAll(cgroupPath, 0555); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		if err := cg.Create(ResourceLimits{
			MemoryLimitBytes: 1,
			IOLimit: LinuxIOLimit{
				Major:    8,
				Minor:    0,
				ReadBPS:  1024,
				WriteBPS: 1024,
			},
		}); err == nil {
			t.Error("Expected error when writing to a read-only cgroup directory")
		}

		// Assert that the cleanup logic in Create() deleted the directory.
		if _, err := os.Stat(cgroupPath); !os.IsNotExist(err) {
			t.Error("Expected cgroup directory to be deleted after partial failure")
		}
	})
}
