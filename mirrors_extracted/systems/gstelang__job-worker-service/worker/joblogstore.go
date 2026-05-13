package worker

import (
	"errors"
	"sync"
)

// ErrJobNotFound is returned when a requested job does not exist.
var ErrJobNotFound = errors.New("job not found")

// JobStatus represents the current state of a job.
type JobStatus int

const (
	StatusInitialized JobStatus = iota
	StatusRunning
	StatusSignaled
	StatusExited
)

// String returns the string representation of a JobStatus.
func (s JobStatus) String() string {
	return [...]string{"Initialized", "Running", "Signaled", "Exited"}[s]
}

// JobLogStore implements both JobLogger and JobDetailsManagement interfaces.
type JobLogStore struct {
	jobs     map[string]JobDetails
	channels map[string]chan []byte
	mu       sync.RWMutex
}

// NewJobLogStore creates and returns a new JobLogStore instance.
func NewJobLogStore() *JobLogStore {
	return &JobLogStore{
		jobs:     make(map[string]JobDetails),
		channels: make(map[string]chan []byte),
	}
}

// AddJob initializes a new job in the store.
func (store *JobLogStore) AddJob(jobID string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.jobs[jobID] = JobDetails{
		JobSummary: JobSummary{Status: StatusInitialized},
	}
}

// AddLog appends a log entry to the specified job.
func (store *JobLogStore) AddLog(jobID string, log []byte) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if job, exists := store.jobs[jobID]; exists {
		job.Logs = append(job.Logs, log)
		store.jobs[jobID] = job
		if ch, ok := store.channels[jobID]; ok {
			ch <- log
		}
	}
}

// GetLogs retrieves all logs for the specified job.
func (store *JobLogStore) GetLogs(jobID string) ([][]byte, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	if job, exists := store.jobs[jobID]; exists {
		return job.Logs, nil
	}
	return nil, ErrJobNotFound
}

// GetOrCreateLogChannel returns an existing log channel for a job or creates a new one.
func (store *JobLogStore) GetOrCreateLogChannel(jobID string) chan []byte {
	store.mu.Lock()
	defer store.mu.Unlock()
	if ch, exists := store.channels[jobID]; exists {
		return ch
	}
	ch := make(chan []byte, 100)
	store.channels[jobID] = ch
	return store.channels[jobID]
}

// GetJobDetails retrieves the full details of a job.
func (store *JobLogStore) GetJobDetails(jobID string) (JobDetails, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	if job, exists := store.jobs[jobID]; exists {
		return job, nil
	}
	return JobDetails{}, ErrJobNotFound
}

// UpdateJobDetails updates the signal and exit code for a job.
func (store *JobLogStore) UpdateJobDetails(jobID string, signal int, exitCode int) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if job, exists := store.jobs[jobID]; exists {
		job.Signal = signal
		job.ExitCode = exitCode
		store.jobs[jobID] = job
	}
}

// AddProcessId sets the process ID for a job.
func (store *JobLogStore) AddProcessId(jobID string, pid int) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if job, exists := store.jobs[jobID]; exists {
		job.PID = pid
		store.jobs[jobID] = job
	}
}

// UpdateJobStatus updates the status of a job.
func (store *JobLogStore) UpdateJobStatus(jobID string, status JobStatus) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if job, exists := store.jobs[jobID]; exists {
		job.Status = status
		store.jobs[jobID] = job
	}
}

// GetJobSummary retrieves a summary of the job details.
func (store *JobLogStore) GetJobSummary(jobID string) (JobSummary, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	if job, exists := store.jobs[jobID]; exists {
		return job.JobSummary, nil
	}
	return JobSummary{}, ErrJobNotFound
}
