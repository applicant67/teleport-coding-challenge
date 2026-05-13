package job

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestNew(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/echo",
		Args: []string{"hello", "world"},
	})

	// Verify initial state
	if j.id == "" {
		t.Error("ID should not be empty")
	}
	if j.exe != "/bin/echo" {
		t.Errorf("Exe() = %q, want %q", j.exe, "/bin/echo")
	}
	if j.State() != Unspecified {
		t.Errorf("State() = %v, want Unspecified", j.State())
	}
	if j.exitCode != -1 {
		t.Errorf("ExitCode() = %d, want -1", j.exitCode)
	}
}

func TestStart_Success(t *testing.T) {
	var stdout bytes.Buffer
	j := New(Config{
		Exe:    "/bin/sh",
		Args:   []string{"-c", "echo hello; sleep 1"},
		Stdout: &stdout,
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Verify running state
	if j.State() != Running {
		t.Errorf("State() = %v, want Running", j.State())
	}

	// Wait for completion
	select {
	case <-j.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("job did not complete in time")
	}

	// Verify exited state
	if j.State() != Exited {
		t.Errorf("State() = %v, want Exited", j.State())
	}
	if j.exitCode != 0 {
		t.Errorf("ExitCode() = %d, want 0", j.exitCode)
	}

	// Verify output was captured
	if got := strings.TrimSpace(stdout.String()); got != "hello" {
		t.Errorf("stdout = %q, want %q", got, "hello")
	}
}

func TestStart_InvalidExecutable(t *testing.T) {
	j := New(Config{Exe: "/nonexistent/path/to/binary"})

	if err := j.Start(); err == nil {
		t.Error("Start() should fail for invalid executable")
	}

	// State should still be unspecified
	if j.State() != Unspecified {
		t.Errorf("State() = %v, want Unspecified", j.State())
	}
}

func TestStart_AlreadyStarted(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/sleep",
		Args: []string{"10"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("first Start() error = %v", err)
	}
	defer j.Stop()

	// Second start should fail
	if err := j.Start(); err == nil {
		t.Error("second Start() should fail")
	}
}

func TestStop_Success(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/sleep",
		Args: []string{"60"}, // Long sleep so we can stop it
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Stop should succeed
	if err := j.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Wait for termination
	select {
	case <-j.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("job did not stop in time")
	}

	// Verify stopped state
	if j.State() != Stopped {
		t.Errorf("State() = %v, want Stopped", j.State())
	}
}

func TestStop_Idempotent(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/sleep",
		Args: []string{"60"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// First stop
	if err := j.Stop(); err != nil {
		t.Fatalf("first Stop() error = %v", err)
	}

	<-j.Done()

	// Second stop should return already stopped
	if err := j.Stop(); !errors.Is(err, ErrAlreadyStopped) {
		t.Errorf("second Stop() = %v, want ErrAlreadyStopped", err)
	}
}

func TestStop_AlreadyStopped(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/sleep",
		Args: []string{"60"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	j.Stop()
	<-j.Done()

	// State is now Stopped
	if err := j.Stop(); !errors.Is(err, ErrAlreadyStopped) {
		t.Errorf("Stop() = %v, want ErrAlreadyStopped", err)
	}
}

func TestStop_AlreadyExited(t *testing.T) {
	j := New(Config{
		Exe: "/bin/true", // Exits immediately
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	<-j.Done()

	// State is now Exited
	if err := j.Stop(); !errors.Is(err, ErrAlreadyExited) {
		t.Errorf("Stop() = %v, want ErrAlreadyExited", err)
	}
}

func TestStop_NotStarted(t *testing.T) {
	j := New(Config{Exe: "/bin/echo"})

	if err := j.Stop(); !errors.Is(err, ErrNotRunning) {
		t.Errorf("Stop() = %v, want ErrNotRunning", err)
	}
}

func TestNonZeroExitCode(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/sh",
		Args: []string{"-c", "exit 42"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	<-j.Done()

	if j.exitCode != 42 {
		t.Errorf("ExitCode() = %d, want 42", j.exitCode)
	}
	if j.State() != Exited {
		t.Errorf("State() = %v, want Exited", j.State())
	}
}

func TestOutputCapture(t *testing.T) {
	var stdout, stderr bytes.Buffer

	j := New(Config{
		Exe:    "/bin/sh",
		Args:   []string{"-c", "echo OUT; echo ERR >&2"},
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	<-j.Done()

	if got := strings.TrimSpace(stdout.String()); got != "OUT" {
		t.Errorf("stdout = %q, want %q", got, "OUT")
	}
	if got := strings.TrimSpace(stderr.String()); got != "ERR" {
		t.Errorf("stderr = %q, want %q", got, "ERR")
	}
}

func TestStatus(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/echo",
		Args: []string{"test"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	<-j.Done()

	status := j.Status()

	if status.ID != j.id {
		t.Errorf("Status.ID = %q, want %q", status.ID, j.id)
	}
	if status.State != Exited {
		t.Errorf("Status.State = %v, want Exited", status.State)
	}
	if status.ExitCode != 0 {
		t.Errorf("Status.ExitCode = %d, want 0", status.ExitCode)
	}
}

func TestDone_ClosesOnExit(t *testing.T) {
	j := New(Config{
		Exe: "/bin/true", // Exits immediately with 0
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Should close quickly
	select {
	case <-j.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("Done() channel was not closed")
	}
}

func TestDone_ClosesOnStop(t *testing.T) {
	j := New(Config{
		Exe:  "/bin/sleep",
		Args: []string{"60"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	j.Stop()

	select {
	case <-j.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("Done() channel was not closed after Stop()")
	}
}

func TestStop_KillsChildProcesses(t *testing.T) {
	// Start a shell that spawns a child process (sleep)
	// The shell itself will exit, but we want to ensure the child is also killed
	j := New(Config{
		Exe:  "/bin/sh",
		Args: []string{"-c", "sleep 60 & sleep 60"},
	})

	if err := j.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Give processes time to spawn
	time.Sleep(500 * time.Millisecond)

	pid := j.cmd.Process.Pid
	if pid <= 0 {
		t.Fatalf("PID() = %d, want > 0", pid)
	}

	pgid, err := unix.Getpgid(pid)
	if err != nil {
		t.Fatalf("Getpgid(%d) error = %v", pid, err)
	}

	// Stop should kill the entire process group
	if err := j.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Wait for termination
	select {
	case <-j.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("job did not stop in time")
	}

	// Process group should be gone.
	// kill(-pgid, 0) checks existence without sending a signal.
	if err := unix.Kill(-pgid, 0); err == nil {
		t.Fatalf("process group %d still exists; expected it to be gone", pgid)
	} else if !errors.Is(err, unix.ESRCH) {
		t.Fatalf("kill(-%d,0) = %v, want ESRCH", pgid, err)
	}

	// Verify state
	if j.State() != Stopped {
		t.Errorf("State() = %v, want Stopped", j.State())
	}
}
