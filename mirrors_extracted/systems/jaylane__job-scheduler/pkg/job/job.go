package job

import (
	"os/exec"
	"slices"
	"time"
)

// Command command used to invoke the job
type Command struct {
	// Name path or program name to execute
	Name string `json:"name"`
	// Args arguments to send to the program/path
	Args []string `json:"args"`
}

// Process details of the process spawned by the job
type Process struct {
	// PID process idententifier
	PID int `json:"pid"`
	// ExitCode exit code of the process
	ExitCode int `json:"exit_code"`
	// Status readable enum of the process status
	Status Status `json:"status"`
}

// Status enum for the status of the job
type Status int

const (
	// TERMINATED status for a job that has completed
	TERMINATED Status = iota
	// RUNNING status for a job that is currently running
	RUNNING
	// STOPPED status for a job that has been stopped
	STOPPED
	// UNKNOWN status for a job that has an unknown status
	UNKNOWN
)

// String implements the Stringer interface for the Status enum
func (s Status) String() string {
	switch s {
	case RUNNING:
		return "running"
	case STOPPED:
		return "stopped"
	case TERMINATED:
		return "terminated"
	default:
		return "unknown"
	}
}

// EnumIndex returns the index of the enum
func (s Status) EnumIndex() int {
	return int(s)
}

// Job details of the job
type Job struct {
	// ID unique identifier for the job
	ID string
	// Cmd command pipeline
	Cmd *exec.Cmd
	// Output output of the job
	Output []byte
	// Process details of the process spawned by the job
	Process *Process
	// StartTime start time of the job
	StartTime time.Time
	// EndTime end time of the job
	EndTime time.Time
}

var nonRunningStatuses = []Status{TERMINATED, STOPPED, UNKNOWN}

// IsRunning returns true if the job is running
func (j *Job) IsRunning() bool {
	return j.Process.ExitCode == 0 && !slices.Contains(nonRunningStatuses, j.Process.Status)
}
