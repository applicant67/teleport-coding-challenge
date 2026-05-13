//go:build linux

package worker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// JobID represents the unique identifier for a job.
type JobID string

// Job is the container for each process initiated by a client.
// It contains the metadata for the Command, Process, and Lifecycle. It also provides access to the Log Stream buffer.
type Job struct {
	mu sync.Mutex

	// ID is a unique UUID generated for the job, distinct from the OS PID. This is the same ID the client will be provided to manage processes.
	ID JobID
	// Specs holds the command from the client.
	Specs ExecutionSpecs
	// Proc holds live handles to the running process.
	Proc ProcessHandle
	// State tracks the current status and timing of the job.
	State Lifecycle
	// Limit defines the resource boundaries for the process on the Linux system.
	Limit ResourceLimits
	// output manages the in-memory buffer for replayable output.
	// It is thread-safe and handles its own locking.
	output *OutputBuffer
	// Cgroup is the resource controller for this job.
	Cgroup *Cgroup
}

// NewJob creates a new Job instance with a generated ID and the given specs.
// cgroupRoot specifies the base directory for cgroups, useful for testing isolation.
func NewJob(specs ExecutionSpecs, cgroupRoot string) (*Job, error) {
	id := JobID(uuid.NewString())
	return &Job{
		ID:    id,
		Specs: specs,
		Proc:  ProcessHandle{},
		State: Lifecycle{
			Status:   StatusUnspecified,
			ExitCode: -1,
		},
		output: NewOutputBuffer(),
		Cgroup: NewCgroup(id, cgroupRoot),
	}, nil
}

// Start initiates the process execution using Ptrace-Controlled Orchestration.
func (j *Job) Start(ctx context.Context) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.State.Status != StatusUnspecified {
		return fmt.Errorf("job already started")
	}

	// Calculate resource limits to get latest I/O device.
	limits, err := getResourceLimits(j.Specs.WorkingDirectory)
	if err != nil {
		return fmt.Errorf("failed to get resource limits: %w", err)
	}
	j.Limit = limits

	if err := j.Cgroup.Create(j.Limit); err != nil {
		j.State.Status = StatusFailed
		return err
	}

	// Create a context with cancellation for the process
	ctx, cancel := context.WithCancel(ctx)
	j.Proc.cancel = cancel
	j.Proc.cmd = j.createCommand(ctx)

	// Ptrace is thread-local. We must lock the goroutine to the current OS thread
	// so that the thread that calls Start() (the tracer) is the same one that calls PtraceDetach().
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := j.Proc.cmd.Start(); err != nil {
		j.abortStart(true)
		return fmt.Errorf("failed to start command: %w", err)
	}

	j.Proc.PID = ProcessID(j.Proc.cmd.Process.Pid)

	exitCode, reaped, err := j.waitForPtraceStop()
	if err != nil {
		j.State.ExitCode = exitCode
		j.abortStart(!reaped)
		return err
	}

	// Move the frozen process into the Cgroup
	if err := j.Cgroup.Attach(j.Proc.PID); err != nil {
		j.abortStart(true)
		return err
	}

	// Detach to allow the process to resume execution
	if err := syscall.PtraceDetach(int(j.Proc.PID)); err != nil {
		j.abortStart(true)
		return fmt.Errorf("failed to detach ptrace: %w", err)
	}

	j.State.Status = StatusRunning
	j.State.StartTime = time.Now().UTC()

	go j.waitAndCleanup()

	return nil
}

// Stop terminates the job using cgroup.kill or SIGKILL.
func (j *Job) Stop() error {
	j.mu.Lock()
	if j.State.Status != StatusRunning {
		j.mu.Unlock()
		return nil
	}

	j.State.Status = StatusStopped
	now := time.Now().UTC()
	j.State.EndTime = &now
	j.mu.Unlock()

	// Try cgroup.kill (Linux 5.14+) to terminate the entire process tree.
	// We explicitly ignore the error because we have a fallback (context cancel) below.
	_ = j.Cgroup.Kill()

	// Cancel the context to ensure the process is killed (fallback) and resources are released.
	if j.Proc.cancel != nil {
		j.Proc.cancel()
	}

	return nil
}

// Stream returns a LogStream iterator that allows clients to read the job's output.
func (j *Job) Stream(ctx context.Context) *LogStream {
	return j.output.Stream(ctx)
}

// Status returns a thread-safe copy of the job's lifecycle state.
func (j *Job) Status() Lifecycle {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.State
}

// createCommand prepares the underlying OS command with the configured context,
// output streams, and system-specific attributes (Ptrace) required for isolation.
func (j *Job) createCommand(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, j.Specs.Command, j.Specs.Arguments...)
	cmd.Dir = j.Specs.WorkingDirectory
	cmd.Stdout = j.output.Stdout()
	cmd.Stderr = j.output.Stderr()

	// Ptrace: true ensures the process pauses before executing the command,
	// allowing us to attach it to the cgroup safely.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	return cmd
}

// waitForPtraceStop blocks until the child process stops at the execve syscall
// (triggered by Ptrace: true). It returns the exit code (if exited), a boolean
// indicating if the process was reaped (exited/signaled), and an error if the
// wait fails or if the process exits prematurely.
func (j *Job) waitForPtraceStop() (int, bool, error) {
	var wstatus syscall.WaitStatus
	// We consume the stop event directly. The final exit event will be generated later
	// when the process actually finishes, so this does not interfere with cmd.Wait().
	if _, err := syscall.Wait4(int(j.Proc.PID), &wstatus, 0, nil); err != nil {
		return -1, false, fmt.Errorf("failed to wait for process stop: %w", err)
	}

	if wstatus.Exited() {
		return wstatus.ExitStatus(), true, fmt.Errorf("process exited prematurely")
	}
	if wstatus.Signaled() {
		return 128 + int(wstatus.Signal()), true, fmt.Errorf("process exited prematurely")
	}

	if !wstatus.Stopped() {
		return -1, false, fmt.Errorf("process did not stop as expected")
	}

	return -1, false, nil
}

// abortStart cleans up resources if the startup sequence fails.
// It kills the process, waits for it to exit (to prevent zombies), deletes the cgroup, and cancels the context.
// killProcess must be false if the process has already been reaped (e.g. by waitForPtraceStop)
// to avoid sending signals to a recycled PID.
func (j *Job) abortStart(killProcess bool) {
	if j.Proc.cmd != nil && j.Proc.cmd.Process != nil {
		// If the process was already reaped (killProcess=false), the PID may have been
		// recycled by the OS. Sending a signal now could kill an unrelated process.
		if killProcess {
			_ = j.Proc.cmd.Process.Kill()
		}
		_ = j.Proc.cmd.Wait() // Reap the process and close internal pipes
	}
	_ = j.Cgroup.Delete()
	if j.Proc.cancel != nil {
		j.Proc.cancel()
	}
	j.State.Status = StatusFailed
	j.output.Close()
}

// waitAndCleanup blocks on the process exit and handles resource cleanup, exit code
// recording, and Cgroup directory removal.
func (j *Job) waitAndCleanup() {
	// Wait for the process to exit
	_ = j.Proc.cmd.Wait()
	j.output.Close()

	exitCode := j.resolveExitCode()
	now := time.Now().UTC()

	j.mu.Lock()
	j.State.ExitCode = exitCode
	j.State.EndTime = &now

	// Only update status if it wasn't manually stopped.
	if j.State.Status != StatusStopped {
		if exitCode == 0 {
			j.State.Status = StatusCompleted
		} else {
			j.State.Status = StatusFailed
		}
	}

	// Ensure context is cancelled to release resources
	if j.Proc.cancel != nil {
		j.Proc.cancel()
	}
	j.mu.Unlock()

	// Ensure any orphaned child processes are terminated so the cgroup can be deleted.
	_ = j.Cgroup.Kill()
	// Cleanup cgroup directory once process is fully reaped
	_ = j.Cgroup.Delete()
}

// resolveExitCode determines the final exit code of the process.
// It handles standard exit codes and reconstructs signal-based terminations (128 + signal).
func (j *Job) resolveExitCode() int {
	if j.Proc.cmd.ProcessState == nil {
		return -1
	}
	code := j.Proc.cmd.ProcessState.ExitCode()
	if code == -1 {
		if status, ok := j.Proc.cmd.ProcessState.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return 128 + int(status.Signal())
		}
	}
	return code
}
