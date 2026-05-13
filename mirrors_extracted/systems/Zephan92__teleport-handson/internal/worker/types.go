package worker

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

// JobStatus represents the current lifecycle state of a worker job.
type JobStatus int32

const (
	// StatusUnspecified indicates the status is unspecified
	StatusUnspecified JobStatus = 0
	// StatusRunning indicates the process has started and is currently active.
	StatusRunning JobStatus = 1
	// StatusStopped indicates the process was manually terminated by a user request.
	StatusStopped JobStatus = 2
	// StatusCompleted indicates the process exited on its own with a zero exit code.
	StatusCompleted JobStatus = 3
	// StatusFailed indicates the process exited with a non-zero exit code or system error.
	StatusFailed JobStatus = 4
)

// ExecutionSpecs defines the parameters provided by the client when a job is first created.
type ExecutionSpecs struct {
	// Command is the path to the executable to be run.
	Command string
	// Arguments is the list of arguments to be passed to the executable.
	Arguments []string
	// WorkingDirectory is the path where the process will execute
	WorkingDirectory string
	// OwnerCommonName is the Common Name extracted from the client's TLS certificate, used for ownership-based authorization.
	OwnerCommonName string
}

// ResourceLimits defines the kernel-level constraints applied to a job.
// These are implemented via Cgroups on Linux.
type ResourceLimits struct {
	// MemoryLimitBytes is the maximum physical memory the process tree can consume.
	MemoryLimitBytes int64
	// CPULimitPercent is the maximum CPU usage allowed (e.g., 0.5 for 50% of a core).
	CPULimitPercent float64
	// IOLimit is the maximum IO Read and Write for the process. It requires specific properties to generate a string necessary for Linux Cgroup limits.
	IOLimit LinuxIOLimit
}

// LinuxIOLimit is the spec needed to generate a structured string for the Linux Cgroup.
type LinuxIOLimit struct {
	// Major is the major device number of the block device to throttle.
	Major uint32
	// Minor is the minor device number of the block device to throttle.
	Minor uint32
	// ReadBPS is the maximum read throughput in bytes per second.
	ReadBPS uint32
	// WriteBPS is the maximum write throughput in bytes per second.
	WriteBPS uint32
}

// String formats the limit into the syntax expected by the cgroup v2 io.max file.
func (l LinuxIOLimit) String() string {
	return fmt.Sprintf("%d:%d rbps=%d wbps=%d", l.Major, l.Minor, l.ReadBPS, l.WriteBPS)
}

// ProcessID represents the operating system's process identifier.
type ProcessID int

// String returns the string representation of the process ID.
func (p ProcessID) String() string {
	return strconv.Itoa(int(p))
}

// ProcessHandle maintains the live OS-level handles and controls for the process.
type ProcessHandle struct {
	// PID is the Operating System's Process ID.
	PID ProcessID
	// cmd is the underlying Go exec.Cmd instance managing the lifecycle.
	cmd *exec.Cmd
	// cancel is the context cancellation function used to signal a stop.
	cancel context.CancelFunc
}

// Lifecycle tracks the temporal state and exit results of a job.
type Lifecycle struct {
	// Status is the current machine-state of the job.
	Status JobStatus
	// ExitCode is the exit code returned by the process upon termination.
	ExitCode int
	// StartTime is the UTC timestamp when the process was successfully started.
	StartTime time.Time
	// EndTime is the UTC timestamp when the process was ended.
	EndTime *time.Time
}

// OutputSource is the source type of each log entry.
type OutputSource int32

const (
	// OutputSourceUnspecified indicates the source is unspecified.
	OutputSourceUnspecified OutputSource = 0
	// OutputSourceStdout is Standard Output from the process.
	OutputSourceStdout OutputSource = 1
	// OutputSourceStderr is Standard Error from the process.
	OutputSourceStderr OutputSource = 2
)

// LogEntry preserves the metadata for a specific chunk of output.
type LogEntry struct {
	// Data contains an unstructured byte array for the current entry from the process.
	Data []byte
	// Source is the type of output log from the process.
	Source OutputSource
	// Timestamp is the timestamp of the Capture event for this log.
	Timestamp time.Time
}
