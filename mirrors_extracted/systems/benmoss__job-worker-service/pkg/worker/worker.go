package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
)

// ErrJobNotFound is returned when a job ID does not exist.
var ErrJobNotFound = errors.New("job not found")

// jobEntry holds the state of a running or completed job.
type jobEntry struct {
	mu       sync.Mutex
	id       uint64
	cmd      *exec.Cmd
	cgroup   *cgroup
	state    JobState
	exitCode int
	command  string
	args     []string
	done     chan struct{} // closed when process exits
	output   *outputWriter
}

// Worker manages the lifecycle of jobs and their associated resources.
type Worker struct {
	config Config
	logger Logger
	mu     sync.RWMutex
	jobs   map[uint64]*jobEntry
	nextID atomic.Uint64
	wg     sync.WaitGroup
}

// New creates a new Worker instance.
// It attempts to clean up cgroup directories from a previous run.
// Returns an error if any cgroups still contain running processes.
func New(config Config) (*Worker, error) {
	if err := cleanupCgroups(config.CgroupRoot); err != nil {
		return nil, fmt.Errorf("cleanup cgroups: %w", err)
	}

	w := &Worker{
		config: config,
		logger: config.Logger,
		jobs:   make(map[uint64]*jobEntry),
	}
	w.nextID.Store(1)

	return w, nil
}

// Start executes a job and returns its unique ID.
func (w *Worker) Start(job Job) (uint64, error) {
	if job.Command == "" {
		return 0, errors.New("command is required")
	}

	id := w.nextID.Add(1) - 1

	var cleanups []func() error
	runCleanups := func() error {
		var errs []error
		for i := len(cleanups) - 1; i >= 0; i-- {
			if err := cleanups[i](); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}

	output, cleanup, err := createLogfile(createLogfileInput{
		LogDir: w.config.LogDir,
		JobID:  id,
	})
	if err != nil {
		return 0, fmt.Errorf("create logfile: %w", err)
	}
	cleanups = append(cleanups, cleanup)

	cg, err := createCgroup(createCgroupInput{
		cgroupRoot:     w.config.CgroupRoot,
		jobID:          id,
		resourceLimits: job.ResourceLimits,
	})
	if err != nil {
		return 0, errors.Join(fmt.Errorf("create cgroup: %w", err), runCleanups())
	}
	cleanups = append(cleanups, cg.remove)

	cgFD, err := cg.fd()
	if err != nil {
		return 0, errors.Join(fmt.Errorf("get cgroup fd: %w", err), runCleanups())
	}
	cleanups = append(cleanups, cgFD.Close)

	cmd := exec.Command(job.Command, job.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		UseCgroupFD: true,
		CgroupFD:    cgFD.FD(),
	}
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Start(); err != nil {
		return 0, errors.Join(fmt.Errorf("start process: %w", err), runCleanups())
	}

	entry := &jobEntry{
		id:      id,
		cmd:     cmd,
		cgroup:  cg,
		state:   JobStateRunning,
		command: job.Command,
		args:    job.Args,
		done:    make(chan struct{}),
		output:  output,
	}

	w.mu.Lock()
	w.jobs[id] = entry
	w.mu.Unlock()

	if err := cgFD.Close(); err != nil && w.logger != nil {
		w.logger.Printf("failed to close cgroup fd for job %d: %v", id, err)
	}

	w.wg.Go(func() { w.waitForProcess(entry) })

	return id, nil
}

// waitForProcess waits for a job's process to exit and updates its state.
func (w *Worker) waitForProcess(entry *jobEntry) {
	err := entry.cmd.Wait()

	// Wait for all child processes in the cgroup to exit.
	entry.cgroup.waitUntilEmpty()

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			entry.exitCode = exitErr.ExitCode()
		} else {
			entry.exitCode = -1
		}
	} else {
		entry.exitCode = 0
	}

	if entry.cmd.ProcessState.Exited() {
		// the job and all subprocesses exited normally
		entry.state = JobStateCompleted
	} else {
		if entry.state == JobStateStopping {
			// the user has signaled for the job to stop
			entry.state = JobStateStopped
		} else {
			// the job was killed but was not signaled to stop, potentially OOM
			entry.state = JobStateKilled
		}
	}
	// broadcast to all log readers to read one last time
	entry.output.markDone()
	close(entry.done)
	if err := entry.output.Close(); err != nil && w.logger != nil {
		w.logger.Printf("failed to close output for job %d: %v", entry.id, err)
	}
}

// Status returns the current state of a job.
func (w *Worker) Status(id uint64) (JobStatus, error) {
	w.mu.RLock()
	entry, ok := w.jobs[id]
	w.mu.RUnlock()

	if !ok {
		return JobStatus{}, ErrJobNotFound
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	return JobStatus{
		ID:       entry.id,
		State:    entry.state,
		ExitCode: entry.exitCode,
		Command:  entry.command,
		Args:     entry.args,
	}, nil
}

// Close shuts down the worker, killing any running jobs, removing job directories, and waiting for them to exit.
func (w *Worker) Close() error {
	w.mu.RLock()
	ids := make([]uint64, 0, len(w.jobs))
	outputs := make([]*outputWriter, 0, len(w.jobs))
	for id, job := range w.jobs {
		ids = append(ids, id)
		outputs = append(outputs, job.output)
	}
	w.mu.RUnlock()

	var errs []error
	for _, id := range ids {
		errs = append(errs, w.Stop(id))
	}

	w.wg.Wait()

	for _, output := range outputs {
		if err := output.RemoveJobDir(); err != nil {
			errs = append(errs, fmt.Errorf("remove job directory %s: %w", filepath.Dir(output.Name()), err))
		}
	}

	return errors.Join(errs...)
}

// Stop terminates a running job and all its child processes.
func (w *Worker) Stop(id uint64) error {
	w.mu.RLock()
	entry, ok := w.jobs[id]
	w.mu.RUnlock()

	if !ok {
		return ErrJobNotFound
	}

	entry.mu.Lock()
	if entry.state == JobStateRunning {
		entry.state = JobStateStopping
		if err := entry.cgroup.kill(); err != nil {
			entry.mu.Unlock()
			return fmt.Errorf("kill job %d: %w", entry.id, err)
		}
	}
	entry.mu.Unlock()
	<-entry.done

	return nil
}

// Logs returns a reader for the job's combined stdout/stderr.
// If follow is true, the reader blocks waiting for new output until the job exits.
// The context controls cancellation; when cancelled, any blocked Read() returns immediately.
func (w *Worker) Logs(ctx context.Context, id uint64, follow bool) (io.ReadCloser, error) {
	w.mu.RLock()
	entry, ok := w.jobs[id]
	w.mu.RUnlock()

	if !ok {
		return nil, ErrJobNotFound
	}

	logfile, err := os.Open(logfilePath(w.config.LogDir, entry.id))
	if err != nil {
		return nil, fmt.Errorf("open logfile for %d: %w", entry.id, err)
	}

	if !follow {
		return logfile, nil
	}

	return &followingReader{
		ctx:    ctx,
		file:   logfile,
		output: entry.output,
	}, nil
}
