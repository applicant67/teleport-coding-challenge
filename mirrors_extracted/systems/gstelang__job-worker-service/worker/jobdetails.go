package worker

// JobSummary represents a summary of a job's details.
type JobSummary struct {
	Status   JobStatus
	PID      int
	ExitCode int
	Signal   int
}

// JobDetails contains job summary details and the logs.
type JobDetails struct {
	JobSummary
	// Need a RingBuffer to hold the most recent entries for scalability since this does not persist to disk.
	// Here's my implementation: https://github.com/gstelang/golang-data-structures/blob/main/other/ringbuffer.go
	Logs [][]byte // store logs in their raw byte form, accommodating any non-UTF-8 characters.
}

// JobDetailsManagement defines operations for managing job details.
type JobDetailsManagement interface {
	GetJobDetails(jobID string) (JobDetails, error)
	UpdateJobDetails(jobID string, signal int, exitCode int)
	AddProcessId(jobID string, pid int)
	UpdateJobStatus(jobID string, status JobStatus)
	GetJobSummary(jobID string) (JobSummary, error)
}
