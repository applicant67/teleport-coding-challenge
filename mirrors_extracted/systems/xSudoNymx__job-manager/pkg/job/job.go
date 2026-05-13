package job

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

var (
	// ErrNotRunning is returned when attempting to stop a job that's not running
	ErrNotRunning = errors.New("job is not running")

	// ErrAlreadyStopped is returned when attempting to stop a job that's already stopped
	ErrAlreadyStopped = errors.New("job has already been stopped")

	// ErrAlreadyExited is returned when attempting to stop a job that has exited
	ErrAlreadyExited = errors.New("job has already exited")

	// ErrUnexpectedState is returned when a job is in an invalid state for the operation
	ErrUnexpectedState = errors.New("job in unexpected state")

	// ErrKillFailed is returned when the process could not be terminated
	ErrKillFailed = errors.New("failed to kill job")

	// ErrAlreadyStarted is returned when attempting to start a job that's already started
	ErrAlreadyStarted = errors.New("job has already been started")

	// ErrStartFailed is returned when the process could not be started
	ErrStartFailed = errors.New("failed to start job")
)

type Job struct {
	// immutable
	id   string
	exe  string
	args []string

	// process writers for output capture
	stdout io.Writer
	stderr io.Writer

	// mutable
	mu        sync.RWMutex
	state     State
	cmd       *exec.Cmd
	startTime time.Time
	endTime   time.Time
	exitCode  int
	waitErr   error

	// Guarantee Stop() logic runs only once
	stopOnce sync.Once

	// done is closed when the job reaches a terminal state
	done chan struct{}
}

// Status represents a snapshot of a job's current state.
type Status struct {
	ID        string
	PID       int
	State     State
	StartedAt time.Time
	EndedAt   time.Time
	ExitCode  int
	WaitErr   error
}

// Config holds options for creating a Job.
type Config struct {
	// Path to executable
	Exe string
	// Arguments to pass the executable
	Args []string
	// Writers for capturing output. if nil, output is discarded
	Stdout io.Writer
	Stderr io.Writer
}

// New creates a new Job with an auto-generated ID.
func New(cfg Config) *Job {
	return NewWithID(uuid.New().String(), cfg)
}

// NewWithID creates a new Job with the specified ID.
func NewWithID(id string, cfg Config) *Job {
	stdout := cfg.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	stderr := cfg.Stderr
	if stderr == nil {
		stderr = io.Discard
	}

	return &Job{
		id:       id,
		exe:      cfg.Exe,
		args:     append([]string(nil), cfg.Args...),
		state:    Unspecified,
		exitCode: -1,
		done:     make(chan struct{}),
		stdout:   stdout,
		stderr:   stderr,
	}
}

func (j *Job) State() State {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.state
}

func (j *Job) Done() <-chan struct{} {
	return j.done
}

// Start begins execution of the job's process.
//
// Returns an error if:
//   - The job was already started
//   - The executable cannot be found or started
func (j *Job) Start() error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.state != Unspecified {
		return fmt.Errorf("%w: %v", ErrAlreadyStarted, j.state)
	}

	cmd := exec.Command(j.exe, j.args...)
	cmd.Stdout = j.stdout
	cmd.Stderr = j.stderr

	// Start process in its own process group.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%w: %v", ErrStartFailed, err)
	}

	// Update state
	j.cmd = cmd
	j.state = Running
	j.startTime = time.Now()

	// Spawn goroutine to wait for process exit
	go j.wait()

	return nil
}

// wait waits for the process to exit and updates state.
func (j *Job) wait() {
	// Wait for process to exit
	err := j.cmd.Wait()

	j.mu.Lock()
	defer j.mu.Unlock()

	j.endTime = time.Now()

	// Extract exit code
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			j.exitCode = exitErr.ExitCode()
		} else {
			j.waitErr = err
		}
	} else {
		j.exitCode = 0
	}

	// Transition to terminal state
	switch j.state {
	case Running:
		j.state = Exited
	case Stopping:
		j.state = Stopped
	default:
		j.waitErr = fmt.Errorf("unexpected state in wait(): %v", j.state)
		j.state = Exited
	}
	close(j.done)
}

// Stop terminates the job's process and all its children.
//
// Returns:
//   - nil if stop was successful
//   - ErrNotRunning if job was never started
//   - ErrAlreadyStopped if job already reached stopped state
//   - ErrAlreadyExited if job already reached exited state
//   - ErrUnexpectedState if job is in an unknown state
//   - ErrKillFailed if process could not be terminated
func (j *Job) Stop() error {
	j.mu.Lock()

	var err error

	switch j.state {
	case Running:
		j.state = Stopping
		pid := j.cmd.Process.Pid
		j.mu.Unlock()

		j.stopOnce.Do(func() {
			// TODO: Handle rare kill errors
			err = j.killJob(pid)
		})
	case Unspecified:
		j.mu.Unlock()
		return ErrNotRunning
	case Stopped, Stopping:
		j.mu.Unlock()
		return ErrAlreadyStopped
	case Exited:
		j.mu.Unlock()
		return ErrAlreadyExited
	default:
		j.mu.Unlock()
		return fmt.Errorf("%w: %v", ErrUnexpectedState, j.state)
	}

	return err
}

func (j *Job) killJob(pid int) error {
	// Kill the entire process group.
	if killErr := unix.Kill(-pid, unix.SIGKILL); killErr != nil {
		if !errors.Is(killErr, unix.ESRCH) {
			return fmt.Errorf("%w: %v", ErrKillFailed, killErr)
		}
	}
	return nil
}

// Status returns a snapshot of the job's current status.
func (j *Job) Status() *Status {
	j.mu.RLock()
	defer j.mu.RUnlock()

	pid := -1

	if j.cmd != nil && j.cmd.Process != nil {
		pid = j.cmd.Process.Pid
	}

	return &Status{
		ID:        j.id,
		PID:       pid,
		State:     j.state,
		StartedAt: j.startTime,
		EndedAt:   j.endTime,
		ExitCode:  j.exitCode,
		WaitErr:   j.waitErr,
	}
}
