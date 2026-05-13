//go:build linux

package worker

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestJob_Lifecycle_Success asserts the full "happy path" of a job:
// it starts, runs to completion, and exits with code 0.
func TestJob_Lifecycle_Success(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Use a command that runs for a bit so we can verify Running state, but exits successfully
	job, err := NewJob(ExecutionSpecs{
		Command:          "sleep",
		Arguments:        []string{"1"},
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(t.Context()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	state := job.Status()
	if state.Status != StatusRunning {
		t.Errorf("Expected status Running, got %v", state.Status)
	}
	if state.ExitCode != -1 {
		t.Errorf("Expected ExitCode -1 while running, got %d", state.ExitCode)
	}
	if state.StartTime.IsZero() {
		t.Error("Expected StartTime to be set")
	}

	// Wait for completion
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job to complete")
		case <-ticker.C:
			s := job.Status()
			if s.Status == StatusCompleted {
				if s.ExitCode != 0 {
					t.Errorf("Expected exit code 0, got %d", s.ExitCode)
				}
				if s.EndTime == nil || s.EndTime.IsZero() {
					t.Error("Expected EndTime to be set")
				}
				return
			}
			if s.Status == StatusFailed {
				t.Fatalf("Job failed unexpectedly with exit code %d", s.ExitCode)
			}
		}
	}
}

// TestJob_Stop asserts that a running job can be terminated successfully.
func TestJob_Stop(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Long running job
	job, err := NewJob(ExecutionSpecs{
		Command:          "sleep",
		Arguments:        []string{"10"},
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(t.Context()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Stop the job
	if err := job.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// The Stop function immediately sets the status to Stopped.
	state := job.Status()
	if state.Status != StatusStopped {
		t.Errorf("Expected status Stopped, got %v", state.Status)
	}
	if state.EndTime == nil || state.EndTime.IsZero() {
		t.Error("Expected EndTime to be set")
	}

	// Wait for the background cleanup goroutine to finish.
	// We verify this by checking if the cgroup directory has been removed.
	cgroupPath := filepath.Join(tmpDir, string(job.ID))
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	isSuccess := false
	for !isSuccess {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job cleanup (cgroup deletion)")
		case <-ticker.C:
			if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
				isSuccess = true
			}
		}
	}

	// Verify ExitCode is set (likely 137 due to SIGKILL from cgroup.kill)
	finalState := job.Status()
	if finalState.ExitCode == -1 {
		t.Error("Expected ExitCode to be set after Stop(), got -1")
	}
}

// TestJob_Start_Errors asserts that the Start method handles various failure modes correctly.
func TestJob_Start_Errors(t *testing.T) {
	t.Parallel()
	t.Run("InvalidCommand", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		// Invalid command
		job, _ := NewJob(ExecutionSpecs{
			Command:          "/bin/does_not_exist",
			WorkingDirectory: tmpDir,
		}, tmpDir)

		if err := job.Start(t.Context()); err == nil {
			t.Error("Expected error starting invalid command")
		}

		state := job.Status()
		if state.Status != StatusFailed {
			t.Errorf("Expected status Failed, got %v", state.Status)
		}
		if state.ExitCode != -1 {
			t.Errorf("Expected ExitCode -1 for failed start, got %d", state.ExitCode)
		}
	})

	t.Run("InvalidWorkingDirectory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		job, _ := NewJob(ExecutionSpecs{
			Command:          "echo",
			WorkingDirectory: "/path/does/not/exist",
		}, tmpDir)

		// Start should fail because getResourceLimits checks the path
		if err := job.Start(t.Context()); err == nil {
			t.Error("Expected error starting with invalid working directory")
		}
	})

	t.Run("CgroupCreateError", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		// Set CgroupRoot to a file to force MkdirAll to fail
		dummyFile := filepath.Join(tmpDir, "dummy_file")
		_ = os.WriteFile(dummyFile, []byte("data"), 0644)

		job, _ := NewJob(ExecutionSpecs{
			Command:          "echo",
			WorkingDirectory: tmpDir,
		}, dummyFile)

		if err := job.Start(t.Context()); err == nil {
			t.Error("Expected error when cgroup creation fails")
		}
	})

	t.Run("ContextCanceled", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		job, _ := NewJob(ExecutionSpecs{
			Command:          "sleep",
			Arguments:        []string{"1"},
			WorkingDirectory: tmpDir,
		}, tmpDir)

		ctx, cancel := context.WithCancel(t.Context())
		cancel() // Cancel immediately

		if err := job.Start(ctx); err == nil {
			t.Error("Expected error starting with canceled context")
		}

		state := job.Status()
		if state.Status != StatusFailed {
			t.Errorf("Expected status Failed, got %v", state.Status)
		}
	})
}

// TestJob_Lifecycle_ExitNonZero asserts that a job that runs but exits with a
// non-zero status code is correctly marked as Failed.
func TestJob_Lifecycle_ExitNonZero(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// Run a command that exits with code 42
	job, err := NewJob(ExecutionSpecs{
		Command:          "sh",
		Arguments:        []string{"-c", "exit 42"},
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for completion
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job to complete")
		case <-ticker.C:
			s := job.Status()
			if s.Status == StatusFailed {
				if s.ExitCode != 42 {
					t.Errorf("Expected exit code 42, got %d", s.ExitCode)
				}
				return
			}
			if s.Status == StatusCompleted {
				t.Fatalf("Job completed successfully but should have failed (exit code %d)", s.ExitCode)
			}
		}
	}
}

// TestJob_SignalTermination asserts that a job killed by an external signal
// (e.g., SIGKILL) is correctly marked as Failed with the proper exit code (128+Signal).
func TestJob_SignalTermination(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	job, err := NewJob(ExecutionSpecs{
		Command:          "sleep",
		Arguments:        []string{"10"},
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(t.Context()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Kill the process externally with SIGKILL (Signal 9)
	pid := int(job.Proc.PID)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		t.Fatalf("Failed to kill process %d: %v", pid, err)
	}

	// Wait for the job to handle the signal and exit
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job to handle SIGKILL")
		case <-ticker.C:
			s := job.Status()
			if s.Status == StatusFailed {
				// Expect 128 + 9 (SIGKILL) = 137
				expectedCode := 128 + int(syscall.SIGKILL)
				if s.ExitCode != expectedCode {
					t.Errorf("Expected exit code %d, got %d", expectedCode, s.ExitCode)
				}
				return
			}
			if s.Status == StatusCompleted {
				t.Fatal("Job completed successfully but should have failed via SIGKILL")
			}
		}
	}
}

// TestJob_ContextCancellation asserts that canceling the parent context of a job
// correctly terminates the process and triggers the cleanup routine.
func TestJob_ContextCancellation(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Use a cancellable context for the job
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	job, err := NewJob(ExecutionSpecs{
		Command:          "sleep",
		Arguments:        []string{"30"},
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Cancel the context, which should trigger the process to be killed.
	cancel()

	// Wait for the cleanup to complete by checking for cgroup deletion.
	cgroupPath := filepath.Join(tmpDir, string(job.ID))
	for i := 0; i < 100; i++ {
		if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
			return // Success
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("Timeout waiting for job cleanup after context cancellation")
}

// TestJob_ImmediateExit asserts that a short-lived process is handled correctly.
func TestJob_ImmediateExit(t *testing.T) {
	t.Parallel()
	if os.Geteuid() != 0 {
		t.Skip("Skipping process tree cleanup test; requires root privileges for cgroups")
	}

	tmpDir := t.TempDir()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	job, err := NewJob(ExecutionSpecs{
		Command:          "true",
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	if err := job.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for completion
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for job to complete")
		case <-ticker.C:
			s := job.Status()
			if s.Status == StatusCompleted {
				if s.ExitCode != 0 {
					t.Errorf("Expected exit code 0, got %d", s.ExitCode)
				}
				return
			}
			if s.Status == StatusFailed {
				t.Fatalf("Job failed unexpectedly")
			}
		}
	}
}

// TestJob_StreamOutput asserts that stdout and stderr are captured and streamed correctly.
func TestJob_StreamOutput(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	// Create a script that prints to stdout and stderr
	scriptPath := filepath.Join(tmpDir, "stream.sh")
	scriptContent := `#!/bin/sh
echo "stdout line"
>&2 echo "stderr line"
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to write script: %v", err)
	}

	job, err := NewJob(ExecutionSpecs{
		Command:          scriptPath,
		WorkingDirectory: tmpDir,
	}, tmpDir)
	if err != nil {
		t.Fatalf("NewJob failed: %v", err)
	}

	// Start streaming before execution
	ctx := t.Context()
	stream := job.Stream(ctx)
	defer stream.Close()

	if err := job.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	var stdout, stderr string
	for {
		entry, err := stream.Next()
		if err == io.EOF {
			break
		}
		if entry.Source == OutputSourceStdout {
			stdout += string(entry.Data)
		} else if entry.Source == OutputSourceStderr {
			stderr += string(entry.Data)
		}
	}

	if !strings.Contains(stdout, "stdout line") {
		t.Errorf("Expected stdout to contain 'stdout line', got %q", stdout)
	}
	if !strings.Contains(stderr, "stderr line") {
		t.Errorf("Expected stderr to contain 'stderr line', got %q", stderr)
	}
}
