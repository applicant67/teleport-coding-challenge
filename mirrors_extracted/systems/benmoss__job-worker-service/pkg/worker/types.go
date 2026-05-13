// Package worker provides a library to manage running Linux processes with
// resource constraints via cgroups v2.
package worker

// Logger defines the logging interface used by the worker.
// Compatible with *log.Logger and fmt.Printf.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Job defines a process to be executed with optional resource constraints.
type Job struct {
	Command        string          // Path to the executable
	Args           []string        // Command-line arguments
	ResourceLimits *ResourceLimits // Optional cgroup limits (nil = no limits)
}

// JobState indicates the current state of a job.
type JobState int

const (
	JobStateRunning   JobState = iota // Process is currently executing
	JobStateCompleted                 // Process exited on its own
	JobStateStopping                  // Stop() called, waiting for exit
	JobStateStopped                   // Terminated by user via Stop()
	JobStateKilled                    // Killed by signal (e.g., OOM) without Stop()
)

// String returns a human-readable representation of the JobState.
func (s JobState) String() string {
	switch s {
	case JobStateRunning:
		return "running"
	case JobStateCompleted:
		return "completed"
	case JobStateStopped:
		return "stopped"
	case JobStateKilled:
		return "killed"
	default:
		return "unknown"
	}
}

// JobStatus represents the current state of a job.
type JobStatus struct {
	ID       uint64   // Unique job identifier
	State    JobState // Current state of the job
	ExitCode int      // Process exit code (only valid if State != JobStateRunning)
	Command  string   // Original command path
	Args     []string // Original arguments
}

// Config contains configuration options for the Worker.
type Config struct {
	// CgroupRoot is the root directory for job cgroups.
	// Defaults to /sys/fs/cgroup/jobworker.
	CgroupRoot string

	// LogDir is the parent directory for job log directories.
	// Defaults to /var/lib/jobworker.
	LogDir string

	// Logger is an optional logger for internal errors.
	// If nil, errors are silently ignored.
	Logger Logger
}

// ResourceLimits specifies cgroups v2 resource constraints for a job.
type ResourceLimits struct {
	CPUMax         float64 // Fraction of CPU (0.5 = half a core, 2.0 = two cores)
	MemoryMaxBytes int64   // Memory limit in bytes
	IOMaxReadBPS   int64   // Max read throughput in bytes/sec
	IOMaxWriteBPS  int64   // Max write throughput in bytes/sec
}
