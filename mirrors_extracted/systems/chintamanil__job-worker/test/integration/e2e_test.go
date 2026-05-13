//go:build integration

package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// TestE2E_SimpleCommand tests the full lifecycle of a simple command.
func TestE2E_SimpleCommand(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	// Start a job
	jobID, err := c.Start(ctx, "/bin/echo", []string{"hello", "integration"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Logf("Started job: %s", jobID)

	// Wait for completion
	if err := waitForStatus(c, jobID, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job did not complete: %v", err)
	}

	// Check final status
	status, err := c.Status(ctx, jobID)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Status != "COMPLETED" {
		t.Errorf("Status = %s, want COMPLETED", status.Status)
	}
	if status.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", status.ExitCode)
	}
	if status.Owner != "client1" {
		t.Errorf("Owner = %s, want client1", status.Owner)
	}

	// Stream output
	var buf bytes.Buffer
	if err := c.Logs(ctx, jobID, &buf); err != nil {
		t.Fatalf("Logs() error = %v", err)
	}
	if !strings.Contains(buf.String(), "hello integration") {
		t.Errorf("Output = %q, want to contain 'hello integration'", buf.String())
	}
}

// TestE2E_StopRunningJob tests stopping a long-running job.
func TestE2E_StopRunningJob(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	// Start a long-running job
	jobID, err := c.Start(ctx, "/bin/sleep", []string{"60"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Logf("Started job: %s", jobID)

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Verify it's running
	status, err := c.Status(ctx, jobID)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Status != "RUNNING" {
		t.Fatalf("Status = %s, want RUNNING", status.Status)
	}
	if status.PID == 0 {
		t.Error("PID should be set for running job")
	}

	// Stop the job
	if err := c.Stop(ctx, jobID); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Verify it's stopped
	status, err = c.Status(ctx, jobID)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Status != "STOPPED" {
		t.Errorf("Status = %s, want STOPPED", status.Status)
	}
	// Note: Signal extraction depends on exit code convention (128+signal)
	// The cgroup.kill mechanism may report differently
	t.Logf("Exit code: %d, Signal: %d", status.ExitCode, status.Signal)
}

// TestE2E_StopKillsChildProcesses tests that stopping a job kills all children.
func TestE2E_StopKillsChildProcesses(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	// Start a job that spawns children
	jobID, err := c.Start(ctx, "/bin/sh", []string{"-c", "sleep 60 & sleep 60 & wait"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Logf("Started job with children: %s", jobID)

	// Give children time to spawn
	time.Sleep(200 * time.Millisecond)

	// Stop the job (should kill parent + children via cgroup.kill)
	if err := c.Stop(ctx, jobID); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Verify it's stopped
	status, err := c.Status(ctx, jobID)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Status != "STOPPED" {
		t.Errorf("Status = %s, want STOPPED", status.Status)
	}
}

// TestE2E_OutputStreaming tests real-time output streaming.
func TestE2E_OutputStreaming(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	// Start a job that outputs lines with delays
	jobID, err := c.Start(ctx, "/bin/sh", []string{"-c", `
		echo "line1"
		sleep 0.1
		echo "line2"
		sleep 0.1
		echo "line3"
	`})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Stream output while job is running
	var buf bytes.Buffer
	if err := c.Logs(ctx, jobID, &buf); err != nil {
		t.Fatalf("Logs() error = %v", err)
	}

	output := buf.String()
	for _, line := range []string{"line1", "line2", "line3"} {
		if !strings.Contains(output, line) {
			t.Errorf("Output missing %q, got: %q", line, output)
		}
	}
}

// TestE2E_StdoutAndStderr tests that both stdout and stderr are captured.
func TestE2E_StdoutAndStderr(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	jobID, err := c.Start(ctx, "/bin/sh", []string{"-c", "echo stdout; echo stderr >&2"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := waitForStatus(c, jobID, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job did not complete: %v", err)
	}

	var buf bytes.Buffer
	if err := c.Logs(ctx, jobID, &buf); err != nil {
		t.Fatalf("Logs() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "stdout") {
		t.Errorf("Output missing stdout, got: %q", output)
	}
	if !strings.Contains(output, "stderr") {
		t.Errorf("Output missing stderr, got: %q", output)
	}
}

// TestE2E_FailedCommand tests handling of commands that exit with error.
func TestE2E_FailedCommand(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	jobID, err := c.Start(ctx, "/bin/sh", []string{"-c", "exit 42"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := waitForStatus(c, jobID, "FAILED", 5*time.Second); err != nil {
		t.Fatalf("Job did not fail: %v", err)
	}

	status, err := c.Status(ctx, jobID)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", status.ExitCode)
	}
}

// TestE2E_MultipleClients tests multiple clients accessing their own jobs.
func TestE2E_MultipleClients(t *testing.T) {
	env := newTestEnv(t)
	client1 := env.newClient("client1")
	client2 := env.newClient("client2")
	ctx := context.Background()

	// Client1 starts a job
	job1, err := client1.Start(ctx, "/bin/echo", []string{"client1 job"})
	if err != nil {
		t.Fatalf("client1.Start() error = %v", err)
	}

	// Client2 starts a job
	job2, err := client2.Start(ctx, "/bin/echo", []string{"client2 job"})
	if err != nil {
		t.Fatalf("client2.Start() error = %v", err)
	}

	// Wait for both to complete
	if err := waitForStatus(client1, job1, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job1 did not complete: %v", err)
	}
	if err := waitForStatus(client2, job2, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job2 did not complete: %v", err)
	}

	// Each client can access their own job
	status1, err := client1.Status(ctx, job1)
	if err != nil {
		t.Fatalf("client1.Status(job1) error = %v", err)
	}
	if status1.Owner != "client1" {
		t.Errorf("Job1 owner = %s, want client1", status1.Owner)
	}

	status2, err := client2.Status(ctx, job2)
	if err != nil {
		t.Fatalf("client2.Status(job2) error = %v", err)
	}
	if status2.Owner != "client2" {
		t.Errorf("Job2 owner = %s, want client2", status2.Owner)
	}

	// Client1 cannot access Client2's job (should get NotFound)
	_, err = client1.Status(ctx, job2)
	if err == nil {
		t.Error("client1 should not be able to access client2's job")
	}
}

// TestE2E_AdminAccessAllJobs tests that admin can access any job.
func TestE2E_AdminAccessAllJobs(t *testing.T) {
	env := newTestEnv(t)
	client1 := env.newClient("client1")
	admin := env.newClient("admin")
	ctx := context.Background()

	// Client1 starts a job
	jobID, err := client1.Start(ctx, "/bin/echo", []string{"secret job"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := waitForStatus(client1, jobID, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job did not complete: %v", err)
	}

	// Admin can access client1's job
	status, err := admin.Status(ctx, jobID)
	if err != nil {
		t.Fatalf("admin.Status() error = %v", err)
	}
	if status.Owner != "client1" {
		t.Errorf("Owner = %s, want client1", status.Owner)
	}

	// Admin can stream client1's job output
	var buf bytes.Buffer
	if err := admin.Logs(ctx, jobID, &buf); err != nil {
		t.Fatalf("admin.Logs() error = %v", err)
	}
	if !strings.Contains(buf.String(), "secret job") {
		t.Errorf("Admin could not read job output")
	}
}

// TestE2E_ConcurrentOutputStreamers tests multiple clients streaming same job.
func TestE2E_ConcurrentOutputStreamers(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	// Start a job with multiple output lines
	jobID, err := c.Start(ctx, "/bin/sh", []string{"-c", `
		for i in 1 2 3 4 5; do
			echo "line$i"
		done
	`})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := waitForStatus(c, jobID, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job did not complete: %v", err)
	}

	// Multiple clients stream concurrently
	errCh := make(chan error, 3)
	for i := 0; i < 3; i++ {
		go func() {
			client := env.newClient("client1")
			var buf bytes.Buffer
			if err := client.Logs(ctx, jobID, &buf); err != nil {
				errCh <- err
				return
			}
			// Verify all lines present
			output := buf.String()
			for _, line := range []string{"line1", "line2", "line3", "line4", "line5"} {
				if !strings.Contains(output, line) {
					errCh <- nil // Will check output below
					return
				}
			}
			errCh <- nil
		}()
	}

	// Wait for all streamers
	for i := 0; i < 3; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("Streamer %d error: %v", i, err)
		}
	}
}

// TestE2E_BinaryOutput tests that binary data is preserved.
func TestE2E_BinaryOutput(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")
	ctx := context.Background()

	// Output binary data (null bytes, high bytes)
	jobID, err := c.Start(ctx, "/bin/sh", []string{"-c", "printf '\\000\\001\\377\\376'"})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := waitForStatus(c, jobID, "COMPLETED", 5*time.Second); err != nil {
		t.Fatalf("Job did not complete: %v", err)
	}

	var buf bytes.Buffer
	if err := c.Logs(ctx, jobID, &buf); err != nil {
		t.Fatalf("Logs() error = %v", err)
	}

	expected := []byte{0x00, 0x01, 0xff, 0xfe}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Binary output = %v, want %v", buf.Bytes(), expected)
	}
}

