package worker

import "context"

// Command represents a system command with arguments
type Command struct {
	Name string
	Args []string
}

// JobController interface (worker) defines the operations on jobs
type JobController interface {
	Start(command Command) (jobID string, err error)
	Stop(jobID string) (ok bool, err error)
	Query(jobID string) (jobSummary JobSummary, err error)
	Stream(ctx context.Context, jobID string) ([][]byte, chan []byte, error)
}
