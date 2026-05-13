package worker

import (
	"errors"
	"os/exec"
	"sync"
	"syscall"

	"github.com/manil/job-worker/pkg/cgroup"
)

// ErrJobNotFound is returned when a job is not found or user is not authorized.
// Using the same error for both cases prevents job ID enumeration attacks.
var ErrJobNotFound = errors.New("job not found")

// ErrJobAlreadyStarted is returned when Start() is called on an already-started job.
var ErrJobAlreadyStarted = errors.New("job already started")

// Sentinel errors for job lifecycle failures.
// These allow callers to distinguish failure types using errors.Is().
var (
	ErrCgroupCreation = errors.New("cgroup creation failed")
	ErrCgroupFDOpen   = errors.New("cgroup fd open failed")
	ErrProcessStart   = errors.New("process start failed")
	ErrJobKill        = errors.New("job kill failed")
)

// JobStatus represents the lifecycle state of a job.
type JobStatus int

const (
	JobStatusUnspecified JobStatus = iota
	JobStatusRunning
	JobStatusCompleted
	JobStatusFailed
	JobStatusStopped
)

// JobState holds all mutable job state. Used both as embedded storage in Job
// and as a snapshot type returned by Job.State() for atomic multi-field reads.
type JobState struct {
	Status     JobStatus
	ExitCode   int
	SignalNum  syscall.Signal // Signal that terminated the process (for STOPPED state)
	Pid        int
	CgroupPath string // Path to job's cgroup directory
	Error      string // Error message if job failed to start
}

// Job represents a single running or completed job.
type Job struct {
	// Immutable fields (set at creation, safe to read without lock)
	ID      string
	Owner   string
	Command string
	Args    []string

	// Resource limits (optional, 0 = unlimited/default)
	CPULimit    float64 // Fraction of one CPU core (0.5 = 50%)
	MemoryLimit int64   // Bytes
	IOWeight    int32   // I/O weight (1-10000, default 100)

	// Mutable process state (protected by mu)
	mu    sync.Mutex
	state JobState

	// Execution state (set during Start, read-only after)
	cmd    *exec.Cmd      // Process handle
	cg     *cgroup.Cgroup // Cgroup for this job
	done   chan struct{}  // Closed when job completes
	Output *MemoryBuffer  // Combined stdout/stderr output buffer
}

// State returns a snapshot of the job's mutable state.
// The returned struct is a copy - safe to read without locks.
func (j *Job) State() JobState {
	j.mu.Lock()
	defer j.mu.Unlock()

	state := j.state
	// Get cgroup path if available
	if j.cg != nil {
		state.CgroupPath = j.cg.Path()
	}
	return state
}

// Cgroup returns the cgroup for this job, or nil if not set.
func (j *Job) Cgroup() *cgroup.Cgroup {
	return j.cg
}

// Cmd returns the exec.Cmd for this job, or nil if not started.
func (j *Job) Cmd() *exec.Cmd {
	return j.cmd
}

// markRunning transitions the job to running state with the given PID.
// Caller must hold j.mu.
func (j *Job) markRunning(pid int) {
	j.state.Status = JobStatusRunning
	j.state.Pid = pid
}

// markCompleted transitions the job to completed state.
// Caller must hold j.mu.
func (j *Job) markCompleted() {
	j.state.Status = JobStatusCompleted
	j.state.ExitCode = 0
}

// markFailed transitions the job to failed state with the given exit code.
// Caller must hold j.mu.
func (j *Job) markFailed(exitCode int, errMsg string) {
	j.state.Status = JobStatusFailed
	j.state.ExitCode = exitCode
	j.state.Error = errMsg
}

// markStopped transitions the job to stopped state (killed by SIGKILL).
// Caller must hold j.mu.
func (j *Job) markStopped(sig syscall.Signal, exitCode int) {
	j.state.Status = JobStatusStopped
	j.state.SignalNum = sig
	j.state.ExitCode = exitCode
}

// markSignaled transitions the job to failed state due to a signal (not SIGKILL).
// Caller must hold j.mu.
func (j *Job) markSignaled(sig syscall.Signal, exitCode int) {
	j.state.Status = JobStatusFailed
	j.state.SignalNum = sig
	j.state.ExitCode = exitCode
}

// JobStore is a thread-safe in-memory store for jobs.
type JobStore struct {
	mu   sync.Mutex
	jobs map[string]*Job
}

// NewJobStore creates a new empty JobStore.
func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*Job),
	}
}

// Add stores a job in the store. If a job with the same ID exists, it is replaced.
func (s *JobStore) Add(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

// Get retrieves a job by ID with ownership verification.
// Returns ErrJobNotFound if the job doesn't exist or user is not authorized.
// Admin users (isAdmin=true) can access any job.
func (s *JobStore) Get(userID, jobID string, isAdmin bool) (*Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return nil, ErrJobNotFound
	}

	// Ownership check: user must own the job, or be admin
	if !isAdmin && job.Owner != userID {
		return nil, ErrJobNotFound // Same error to prevent enumeration
	}

	return job, nil
}

// GetByID retrieves a job by ID without ownership check.
// Use this only for internal operations where authorization has already been verified.
func (s *JobStore) GetByID(jobID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[jobID]
	return job, ok
}

// Delete removes a job from the store. No-op if job doesn't exist.
func (s *JobStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
}

// Count returns the number of jobs in the store.
func (s *JobStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.jobs)
}
