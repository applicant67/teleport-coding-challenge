package worker_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/benmoss/job-worker-service/pkg/worker"
)

// testCgroupRoot returns a cgroup root path for tests.
func testCgroupRoot(t *testing.T) string {
	t.Helper()

	cgroupRoot := fmt.Sprintf("/sys/fs/cgroup/jobworker-test-%d", os.Getpid())

	if err := os.MkdirAll(cgroupRoot, 0o755); err != nil {
		t.Fatalf("cannot create cgroup directory (need root or cgroup delegation): %v", err)
	}

	t.Cleanup(func() {
		entries, _ := os.ReadDir(cgroupRoot)
		for _, e := range entries {
			if e.IsDir() {
				if err := os.Remove(filepath.Join(cgroupRoot, e.Name())); err != nil {
					t.Logf("failed to remove cgroup entry: %v", err)
				}
			}
		}
		if err := os.Remove(cgroupRoot); err != nil {
			t.Logf("failed to remove cgroup: %v", err)
		}
	})

	return cgroupRoot
}

// waitForJobStateInput contains parameters for waitForJobState.
type waitForJobStateInput struct {
	Worker      *worker.Worker
	JobID       uint64
	Timeout     time.Duration
	PollInterval time.Duration
}

// waitForJobState polls the job status until it's no longer running or timeout occurs.
// Returns the final status and the duration it took.
func waitForJobState(t *testing.T, input waitForJobStateInput) (worker.JobStatus, time.Duration) {
	t.Helper()

	if input.PollInterval == 0 {
		input.PollInterval = 10 * time.Millisecond
	}

	start := time.Now()
	deadline := start.Add(input.Timeout)
	var status worker.JobStatus
	var err error

	for time.Now().Before(deadline) {
		status, err = input.Worker.Status(input.JobID)
		if err != nil {
			t.Fatalf("Status() error: %v", err)
		}
		if status.State != worker.JobStateRunning {
			break
		}
		time.Sleep(input.PollInterval)
	}

	return status, time.Since(start)
}

func TestWorker(t *testing.T) {
	cgroupRoot := testCgroupRoot(t)

	cfg := worker.Config{
		CgroupRoot: cgroupRoot,
	}
	w, err := worker.New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	t.Cleanup(func() {
		if err := w.Close(); err != nil {
			t.Logf("failed to close worker: %v", err)
		}
	})

	t.Run("start job and check running status", func(t *testing.T) {
		job := worker.Job{
			Command: "/bin/sleep",
			Args:    []string{"10"},
		}
		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		status, err := w.Status(id)
		if err != nil {
			t.Fatalf("Status() error: %v", err)
		}
		if status.State != worker.JobStateRunning {
			t.Errorf("State = %v, want %v", status.State, worker.JobStateRunning)
		}
		if status.ID != id {
			t.Errorf("ID = %d, want %d", status.ID, id)
		}
		if status.Command != job.Command {
			t.Errorf("Command = %q, want %q", status.Command, job.Command)
		}
	})

	t.Run("job that exits has completed status", func(t *testing.T) {
		job := worker.Job{
			Command: "/bin/true",
		}
		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		status, _ := waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id,
			Timeout: 5 * time.Second,
		})

		if status.State != worker.JobStateCompleted {
			t.Errorf("State = %v, want %v", status.State, worker.JobStateCompleted)
		}
		if status.ExitCode != 0 {
			t.Errorf("ExitCode = %d, want 0", status.ExitCode)
		}
	})

	t.Run("job that fails has completed status with non-zero exit code", func(t *testing.T) {
		job := worker.Job{
			Command: "/bin/false",
		}
		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		status, _ := waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id,
			Timeout: 5 * time.Second,
		})

		if status.State != worker.JobStateCompleted {
			t.Errorf("State = %v, want %v", status.State, worker.JobStateCompleted)
		}
		if status.ExitCode == 0 {
			t.Error("ExitCode = 0, want non-zero")
		}
	})

	t.Run("status of nonexistent job returns error", func(t *testing.T) {
		_, err := w.Status(999999)
		if err == nil {
			t.Fatal("Status() expected error for nonexistent job")
		}
	})

	t.Run("start with empty command returns error", func(t *testing.T) {
		job := worker.Job{}
		_, err := w.Start(job)
		if err == nil {
			t.Fatal("Start() expected error for empty command")
		}
	})

	t.Run("stop terminates a running job", func(t *testing.T) {
		job := worker.Job{
			Command: "/bin/sleep",
			Args:    []string{"infinity"},
		}

		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		if err := w.Stop(id); err != nil {
			t.Fatalf("Stop() error: %v", err)
		}

		status, err := w.Status(id)
		if err != nil {
			t.Fatalf("Status() error: %v", err)
		}

		if status.State != worker.JobStateStopped {
			t.Errorf("State = %v, want %v", status.State, worker.JobStateKilled)
		}
	})

	t.Run("logs of a running job", func(t *testing.T) {
		job := worker.Job{
			Command: "/bin/sh",
			Args:    []string{"-c", "echo hello && echo world"},
		}

		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		// Wait for job to complete to ensure logs are fully written
		waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id,
			Timeout: 3 * time.Second,
		})

		output, err := w.Logs(context.Background(), id, false)
		if err != nil {
			t.Fatalf("Logs() error: %v", err)
		}
		defer output.Close()

		logs, err := io.ReadAll(output)
		if err != nil {
			t.Fatalf("ReadAll() error: %v", err)
		}

		expected := "hello\nworld\n"
		if string(logs) != expected {
			t.Errorf("logs = %q, want %q", string(logs), expected)
		}
	})

	t.Run("follow continuously streams data from long-running job", func(t *testing.T) {
		job := worker.Job{
			Command: "/bin/sh",
			Args:    []string{"-c", "while true; do echo data; sleep 0.1; done"},
		}

		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}
		defer w.Stop(id)

		output, err := w.Logs(context.Background(), id, true)
		if err != nil {
			t.Fatalf("Logs() error: %v", err)
		}
		defer output.Close()

		// Read multiple chunks to verify streaming continues
		buf := make([]byte, 1024)
		totalBytes := 0
		readsWithData := 0
		deadline := time.Now().Add(1 * time.Second)

		for time.Now().Before(deadline) {
			n, err := output.Read(buf)
			if n > 0 {
				totalBytes += n
				readsWithData++
			}
			if err == io.EOF {
				t.Fatal("unexpected EOF - stream should continue while job is running")
			}
			if err != nil {
				t.Fatalf("Read() error: %v", err)
			}
		}

		// Verify we received data continuously
		if totalBytes == 0 {
			t.Error("expected to receive data from streaming job, got 0 bytes")
		}
		if readsWithData < 3 {
			t.Errorf("expected multiple reads with data, got %d", readsWithData)
		}
	})

	t.Run("memory limit OOM kills job", func(t *testing.T) {
		// Create a job that allocates memory until OOM killed
		// Use Python to actually allocate and hold memory in a list
		job := worker.Job{
			Command: "/usr/bin/python3",
			Args: []string{"-c", `
import time
data = []
try:
    while True:
        data.append(' ' * (1024 * 1024))  # Allocate 1MB
        time.sleep(0.01)
except:
    pass
`},
			ResourceLimits: &worker.ResourceLimits{
				MemoryMaxBytes: 10 * 1024 * 1024, // 10MB limit
			},
		}

		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		status, _ := waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id,
			Timeout: 5 * time.Second,
		})

		if status.State == worker.JobStateRunning {
			t.Fatal("job should have been killed by OOM, but is still running")
		}

		// Job should be killed (OOM kill shows as SIGKILL)
		if status.State != worker.JobStateKilled {
			t.Errorf("State = %v, want %v (OOM killed)", status.State, worker.JobStateKilled)
		}
	})

	t.Run("cpu limit throttles job", func(t *testing.T) {
		// Run the same CPU-bound workload twice: once without limits, once with 20% CPU limit
		// The throttled job should take significantly longer
		workload := []string{"-c", "x=0; i=0; while [ $i -lt 200000 ]; do x=$((x + i)); i=$((i + 1)); done"}

		// Run without limits
		jobUnlimited := worker.Job{
			Command:        "/bin/sh",
			Args:           workload,
			ResourceLimits: nil,
		}

		id1, err := w.Start(jobUnlimited)
		if err != nil {
			t.Fatalf("Start() unlimited job error: %v", err)
		}

		statusUnlimited, unlimitedDuration := waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id1,
			Timeout: 10 * time.Second,
		})

		if statusUnlimited.State != worker.JobStateCompleted {
			t.Fatalf("Unlimited job state = %v, want %v", statusUnlimited.State, worker.JobStateCompleted)
		}

		// Run with 20% CPU limit
		jobLimited := worker.Job{
			Command: "/bin/sh",
			Args:    workload,
			ResourceLimits: &worker.ResourceLimits{
				CPUMax: 0.2, // 20% of one CPU
			},
		}

		id2, err := w.Start(jobLimited)
		if err != nil {
			t.Fatalf("Start() limited job error: %v", err)
		}

		statusLimited, limitedDuration := waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id2,
			Timeout: 30 * time.Second,
		})

		if statusLimited.State != worker.JobStateCompleted {
			t.Fatalf("Limited job state = %v, want %v", statusLimited.State, worker.JobStateCompleted)
		}

		t.Logf("Unlimited job: %v, Limited (20%% CPU): %v, Ratio: %.2fx",
			unlimitedDuration, limitedDuration, float64(limitedDuration)/float64(unlimitedDuration))

		// The throttled job should take at least 2x as long
		// With 20% CPU, it should theoretically take 5x longer, but allow for scheduling overhead
		if limitedDuration < unlimitedDuration*2 {
			t.Errorf("CPU throttling not effective: unlimited took %v, limited (20%%) took %v (expected at least 2x slower)",
				unlimitedDuration, limitedDuration)
		}
	})

	t.Run("job without resource limits runs successfully", func(t *testing.T) {
		job := worker.Job{
			Command:        "/bin/echo",
			Args:           []string{"no limits"},
			ResourceLimits: nil, // No limits
		}

		id, err := w.Start(job)
		if err != nil {
			t.Fatalf("Start() error: %v", err)
		}

		status, _ := waitForJobState(t, waitForJobStateInput{
			Worker:  w,
			JobID:   id,
			Timeout: 3 * time.Second,
		})

		if status.State != worker.JobStateCompleted {
			t.Errorf("State = %v, want %v", status.State, worker.JobStateCompleted)
		}
		if status.ExitCode != 0 {
			t.Errorf("ExitCode = %d, want 0", status.ExitCode)
		}
	})
}
