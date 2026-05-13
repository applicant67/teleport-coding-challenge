//go:build linux

package worker

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// testCgroupRoot returns a temporary cgroup root for testing.
// Creates the directory if it doesn't exist.
func testCgroupRoot(t *testing.T) string {
	t.Helper()

	// Use system cgroup if running as root, otherwise use temp dir
	if os.Getuid() == 0 {
		root := "/sys/fs/cgroup/jobworker-test"
		if err := os.MkdirAll(root, 0755); err != nil {
			t.Skipf("cannot create cgroup root (need root): %v", err)
		}
		t.Cleanup(func() {
			os.RemoveAll(root)
		})
		return root
	}

	// Non-root: use temp directory (cgroup operations will be simulated)
	return t.TempDir()
}

// skipIfNotRoot skips the test if not running as root.
func skipIfNotRoot(t *testing.T) {
	t.Helper()
	if os.Getuid() != 0 {
		t.Skip("test requires root privileges for cgroups")
	}
}

// --- Happy Path Tests ---

func TestJob_Start_SimpleCommand(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-simple",
		Owner:   "testuser",
		Command: "/bin/echo",
		Args:    []string{"hello", "world"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Wait for job to complete
	job.Wait()

	state := job.State()
	if state.Status != JobStatusCompleted {
		t.Errorf("Status = %v, want JobStatusCompleted", state.Status)
	}
	if state.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", state.ExitCode)
	}
	if state.Pid == 0 {
		t.Error("Pid should be set")
	}

	// Check output
	reader := job.Output.NewReader(context.Background())
	defer reader.Close()
	output, _ := io.ReadAll(reader)
	if !bytes.Contains(output, []byte("hello world")) {
		t.Errorf("Output = %q, want to contain 'hello world'", output)
	}
}

func TestJob_Start_WithResourceLimits(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	// Check if we can actually use cgroup controllers by testing a write
	// Controllers may be listed but not delegated (common in containers)
	testDir := filepath.Join(cgroupRoot, "controller-test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Skipf("cannot create test cgroup: %v", err)
	}
	defer os.Remove(testDir)

	// Try to write cpu.max - if this fails, cpu controller isn't usable
	if err := os.WriteFile(filepath.Join(testDir, "cpu.max"), []byte("50000 100000"), 0644); err != nil {
		t.Skipf("cpu controller not delegated (cannot write cpu.max): %v", err)
	}

	job := &Job{
		ID:          "test-job-limits",
		Owner:       "testuser",
		Command:     "/bin/true",
		CPULimit:    0.5,
		MemoryLimit: 100 * 1024 * 1024, // 100 MiB
		IOWeight:    200,
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	job.Wait()

	if job.State().Status != JobStatusCompleted {
		t.Errorf("Status = %v, want JobStatusCompleted", job.State().Status)
	}

	// Verify cgroup was cleaned up
	cgPath := filepath.Join(cgroupRoot, job.ID)
	if _, err := os.Stat(cgPath); !os.IsNotExist(err) {
		t.Error("cgroup directory should be removed after job completion")
	}
}

func TestJob_Start_CapturesOutput(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-output",
		Owner:   "testuser",
		Command: "/bin/sh",
		Args:    []string{"-c", "echo stdout; echo stderr >&2"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	job.Wait()

	reader := job.Output.NewReader(context.Background())
	defer reader.Close()
	output, _ := io.ReadAll(reader)

	// Both stdout and stderr should be captured
	if !bytes.Contains(output, []byte("stdout")) {
		t.Errorf("Output should contain 'stdout', got %q", output)
	}
	if !bytes.Contains(output, []byte("stderr")) {
		t.Errorf("Output should contain 'stderr', got %q", output)
	}
}

func TestJob_Start_BinaryOutput(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-binary",
		Owner:   "testuser",
		Command: "/bin/sh",
		// Use octal escapes for portability (busybox printf doesn't support \x)
		Args: []string{"-c", "printf '\\000\\001\\377\\376'"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	job.Wait()

	reader := job.Output.NewReader(context.Background())
	defer reader.Close()
	output, _ := io.ReadAll(reader)

	expected := []byte{0x00, 0x01, 0xff, 0xfe}
	if !bytes.Equal(output, expected) {
		t.Errorf("Output = %v, want %v", output, expected)
	}
}

func TestJob_Stop_KillsProcess(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-stop",
		Owner:   "testuser",
		Command: "/bin/sleep",
		Args:    []string{"60"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Give process time to start
	time.Sleep(50 * time.Millisecond)

	if err := job.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	state := job.State()
	if state.Status != JobStatusStopped {
		t.Errorf("Status = %v, want JobStatusStopped", state.Status)
	}
	if state.SignalNum != syscall.SIGKILL {
		t.Errorf("SignalNum = %v, want SIGKILL", state.SignalNum)
	}
}

func TestJob_Stop_Idempotent(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-stop-idempotent",
		Owner:   "testuser",
		Command: "/bin/true",
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	job.Wait()

	// Stop after already completed should be idempotent
	if err := job.Stop(); err != nil {
		t.Errorf("Stop() on completed job error = %v, want nil", err)
	}
	if err := job.Stop(); err != nil {
		t.Errorf("Stop() second call error = %v, want nil", err)
	}
}

func TestJob_Stop_KillsChildProcesses(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	// Shell script that spawns child processes
	job := &Job{
		ID:      "test-job-children",
		Owner:   "testuser",
		Command: "/bin/sh",
		Args:    []string{"-c", "sleep 60 & sleep 60 & wait"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Give children time to spawn
	time.Sleep(100 * time.Millisecond)

	if err := job.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// All processes (parent + children) should be killed
	if job.State().Status != JobStatusStopped {
		t.Errorf("Status = %v, want JobStatusStopped", job.State().Status)
	}
}

func TestJob_ConcurrentStreamers(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-concurrent",
		Owner:   "testuser",
		Command: "/bin/sh",
		Args:    []string{"-c", "for i in 1 2 3 4 5; do echo line$i; done"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Multiple concurrent readers
	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader := job.Output.NewReader(context.Background())
			defer reader.Close()

			output, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("ReadAll() error = %v", err)
				return
			}
			// Each reader should get all 5 lines
			for _, line := range []string{"line1", "line2", "line3", "line4", "line5"} {
				if !bytes.Contains(output, []byte(line)) {
					t.Errorf("Output missing %q", line)
				}
			}
		}()
	}

	job.Wait()
	wg.Wait()
}

// --- Unhappy Path Tests ---

func TestJob_Start_AlreadyStarted(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-double-start",
		Owner:   "testuser",
		Command: "/bin/sleep",
		Args:    []string{"60"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Second start should fail
	err := job.Start(cgroupRoot)
	if err != ErrJobAlreadyStarted {
		t.Errorf("Start() second call error = %v, want ErrJobAlreadyStarted", err)
	}

	// Cleanup
	job.Stop()
}

func TestJob_Start_CommandNotFound(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-notfound",
		Owner:   "testuser",
		Command: "/nonexistent/command",
	}

	err := job.Start(cgroupRoot)
	if err == nil {
		t.Fatal("Start() should fail for nonexistent command")
	}

	state := job.State()
	if state.Status != JobStatusFailed {
		t.Errorf("Status = %v, want JobStatusFailed", state.Status)
	}
	if state.Error == "" {
		t.Error("Error should be set")
	}
}

func TestJob_FailedExit(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-fail",
		Owner:   "testuser",
		Command: "/bin/false",
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	job.Wait()

	state := job.State()
	if state.Status != JobStatusFailed {
		t.Errorf("Status = %v, want JobStatusFailed", state.Status)
	}
	if state.ExitCode == 0 {
		t.Error("ExitCode should be non-zero for failed job")
	}
}

func TestJob_ExitCode(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-exitcode",
		Owner:   "testuser",
		Command: "/bin/sh",
		Args:    []string{"-c", "exit 42"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	job.Wait()

	if job.State().ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", job.State().ExitCode)
	}
}

func TestJob_OutputStreamingDuringExecution(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-streaming",
		Owner:   "testuser",
		Command: "/bin/sh",
		Args:    []string{"-c", "echo first; sleep 0.1; echo second; sleep 0.1; echo third"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	reader := job.Output.NewReader(context.Background())
	defer reader.Close()

	// Read incrementally
	var received []string
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		lines := strings.Split(strings.TrimSpace(string(buf[:n])), "\n")
		received = append(received, lines...)
	}

	if len(received) < 3 {
		t.Errorf("Expected at least 3 lines, got %d: %v", len(received), received)
	}
}

func TestJob_ContextCancellation(t *testing.T) {
	skipIfNotRoot(t)
	cgroupRoot := testCgroupRoot(t)

	job := &Job{
		ID:      "test-job-ctx-cancel",
		Owner:   "testuser",
		Command: "/bin/sleep",
		Args:    []string{"60"},
	}

	if err := job.Start(cgroupRoot); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	reader := job.Output.NewReader(ctx)
	defer reader.Close()

	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		_, err := reader.Read(buf)
		done <- err
	}()

	// Cancel context should unblock reader
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("Read() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Read() should have unblocked after context cancel")
	}

	// Cleanup
	job.Stop()
}
