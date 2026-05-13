//go:build linux

package worker

import (
	"testing"
)

// TestNewJob asserts the successful creation of a Job and its error handling.
func TestNewJob(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		specs := ExecutionSpecs{
			Command:          "echo",
			Arguments:        []string{"hello"},
			WorkingDirectory: tmpDir,
		}

		job, err := NewJob(specs, tmpDir)
		if err != nil {
			t.Fatalf("NewJob failed: %v", err)
		}

		if job.ID == "" {
			t.Error("Job ID should not be empty")
		}
		if job.State.Status != StatusUnspecified {
			t.Errorf("Expected status Unspecified, got %v", job.State.Status)
		}
		if job.State.ExitCode != -1 {
			t.Errorf("Expected initial ExitCode -1, got %d", job.State.ExitCode)
		}
		if job.output == nil {
			t.Error("OutputBuffer should be initialized")
		}
		if job.Cgroup == nil {
			t.Error("Cgroup should be initialized")
		}
	})
}

// TestJob_Stop_NotRunning asserts that calling Stop on a job that is not
// currently running is a safe no-op.
func TestJob_Stop_NotRunning(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	job, _ := NewJob(ExecutionSpecs{
		Command:          "echo",
		WorkingDirectory: tmpDir,
	}, tmpDir)

	if err := job.Stop(); err != nil {
		t.Errorf("Expected no error stopping unstarted job, got: %v", err)
	}
}

// TestJob_Start_AlreadyStarted asserts that attempting to start a job that is
// already in the Running state returns an error.
func TestJob_Start_AlreadyStarted(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	job, _ := NewJob(ExecutionSpecs{Command: "echo", WorkingDirectory: tmpDir}, tmpDir)
	job.State.Status = StatusRunning

	if err := job.Start(t.Context()); err == nil {
		t.Error("Expected error when starting already running job")
	}
}
