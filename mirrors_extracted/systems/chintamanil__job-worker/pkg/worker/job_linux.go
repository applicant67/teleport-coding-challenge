//go:build linux

package worker

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/manil/job-worker/pkg/cgroup"
)

// Start launches the job's command in a new cgroup with resource limits.
// The process is placed in its own process group and cgroup atomically at fork.
// Returns ErrJobAlreadyStarted if the job has already been started.
func (j *Job) Start(cgroupRoot string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.state.Status != JobStatusUnspecified {
		return ErrJobAlreadyStarted
	}

	// Initialize output buffer
	j.Output = NewMemoryBuffer()
	j.done = make(chan struct{})

	// Create cgroup with resource limits
	j.cg = cgroup.New(cgroupRoot, j.ID)
	limits := cgroup.Limits{
		CPUQuota:    j.CPULimit,
		MemoryBytes: j.MemoryLimit,
		IOWeight:    j.IOWeight,
	}

	if err := j.cg.Create(limits); err != nil {
		j.markFailed(-1, fmt.Sprintf("failed to create cgroup: %v", err))
		close(j.done)
		return fmt.Errorf("%w: %w", ErrCgroupCreation, err)
	}

	// Open cgroup FD for atomic process assignment
	cgroupFD, err := j.cg.OpenFD()
	if err != nil {
		j.cleanup()
		j.markFailed(-1, fmt.Sprintf("failed to open cgroup fd: %v", err))
		close(j.done)
		return fmt.Errorf("%w: %w", ErrCgroupFDOpen, err)
	}

	// Create command - no shell interpretation
	j.cmd = exec.Command(j.Command, j.Args...)

	// Process isolation:
	// - Setpgid: new process group (prevents signal bleed from server)
	// - UseCgroupFD: atomic cgroup assignment at fork (prevents child escape)
	j.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:     true,
		UseCgroupFD: true,
		CgroupFD:    cgroupFD,
	}

	// Capture combined stdout/stderr to preserve output ordering
	j.cmd.Stdout = j.Output
	j.cmd.Stderr = j.Output

	// Start the process
	if err := j.cmd.Start(); err != nil {
		j.cg.CloseFD()
		j.cleanup()
		j.markFailed(-1, fmt.Sprintf("failed to start process: %v", err))
		j.Output.Close()
		close(j.done)
		return fmt.Errorf("%w: %w", ErrProcessStart, err)
	}

	// Close cgroup FD after start - process is already in cgroup
	j.cg.CloseFD()

	j.markRunning(j.cmd.Process.Pid)

	// Launch goroutine to wait for process exit
	go j.wait()

	return nil
}

// Stop terminates the job by killing all processes in its cgroup.
// Uses cgroup.kill for atomic termination of the entire process tree.
// Returns nil if the job is already stopped (idempotent).
func (j *Job) Stop() error {
	j.mu.Lock()

	if j.state.Status != JobStatusRunning {
		j.mu.Unlock()
		// Idempotent: already stopped/completed/failed
		return nil
	}

	j.mu.Unlock()

	// Kill entire process tree via cgroup
	if err := j.cg.Kill(); err != nil {
		return fmt.Errorf("%w: %w", ErrJobKill, err)
	}

	// Wait for the job to finish (wait goroutine will update state)
	<-j.done

	return nil
}

// Wait blocks until the job completes (either naturally or via Stop).
func (j *Job) Wait() {
	<-j.done
}

// Done returns a channel that's closed when the job completes.
func (j *Job) Done() <-chan struct{} {
	return j.done
}

// wait monitors process exit and updates job state.
// Called as a goroutine from Start().
func (j *Job) wait() {
	// Block until process exits
	err := j.cmd.Wait()

	j.mu.Lock()
	defer j.mu.Unlock()

	// Close output buffer to signal EOF to readers
	j.Output.Close()

	// Determine final status based on exit
	if err == nil {
		// Clean exit (exit code 0)
		j.markCompleted()
	} else {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()

			// Check if killed by signal
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					sig := status.Signal()
					// SIGKILL from cgroup.kill indicates stopped by user
					if sig == syscall.SIGKILL {
						j.markStopped(sig, exitCode)
					} else {
						j.markSignaled(sig, exitCode)
					}
				} else {
					j.markFailed(exitCode, "")
				}
			} else {
				j.markFailed(exitCode, "")
			}
		} else {
			// Other error (shouldn't happen after successful Start)
			j.markFailed(-1, err.Error())
		}
	}

	// Cleanup cgroup
	j.cleanup()

	// Signal completion
	close(j.done)
}

// cleanup removes the cgroup directory.
// Must be called after all processes have exited.
func (j *Job) cleanup() {
	if j.cg != nil {
		j.cg.Remove()
	}
}
