# Job worker service
## Overview
This document specifies the design of a Go library to manage running Linux processes, a gRPC service that uses this library, and a CLI for interacting with the service. The library and service support running arbitrary Linux processes within the constraints offered by Linux cgroups and offer real-time streaming of process output.
## CLI Usage
```shell-session
$ jobctl start --memory 512M --cpu 0.5 -- /usr/bin/python script.py
Job 1 created

$ jobctl status 1
Status:    running
Command:   /usr/bin/python
Arguments: script.py

$ jobctl logs 1 # displays current output and exits
$ jobctl logs --follow 1 # streams until job exits or ctrl-c

$ jobctl stop 1
Job 1 stopped

$ jobctl status 1
Status:    exited
Exit code: 0
Command:   /usr/bin/python
Arguments: script.py

$ jobctl status 999
Error: job not found
```

## gRPC API

```protobuf
// ResourceLimits specifies cgroups v2 resource constraints.
// All fields are optional; zero values mean no limit.
message ResourceLimits {
	double cpu_max = 1;          // Fraction of CPU (0.5 = half core, 2.0 = two cores)
	int64 memory_max_bytes = 2;  // Memory limit in bytes
	int64 io_max_read_bps = 3;   // Max read throughput in bytes/sec
	int64 io_max_write_bps = 4;  // Max write throughput in bytes/sec
}

// StartJobRequest specifies a process to execute.
message StartJobRequest {
	string command = 1;           // Path to executable (required)
	repeated string args = 2;     // Command-line arguments
	ResourceLimits limits = 3;    // Resource limits (optional)
}

message StartJobResponse {
	uint64 job_id = 1;  // Unique identifier for the created job
}

message StopJobRequest {
	uint64 job_id = 1;
}

message StopJobResponse {}

message GetStatusRequest {
	uint64 job_id = 1;
}

// JobState indicates the current state of a job.
enum JobState {
	RUNNING = 0;    // Process is currently executing
	COMPLETED = 1;  // Process exited on its own
	STOPPING = 2;   // StopJob called, waiting on exit
	STOPPED = 3;    // Terminated by user via StopJob
	KILLED = 4;     // Killed by signal (e.g., OOM) without StopJob
}

// GetStatusResponse returns the current state of a job.
message GetStatusResponse {
	uint64 job_id = 1;
	JobState state = 2;
	int32 exit_code = 3;  // Exit code (only valid if state != RUNNING)
	string command = 4;
	repeated string args = 5;
}

message GetLogsRequest {
	uint64 job_id = 1;
	bool follow = 2;  // If true, stream waits for new output until job exits
}

// GetLogsResponse contains a chunk of output data.
// Output is opaque binary data with no encoding assumptions.
message GetLogsResponse {
	bytes data = 1;
}

// JobWorker manages Linux processes with resource constraints.
service JobWorker {
	// StartJob creates and executes a new job. Returns after the process starts.
	// Errors: INVALID_ARGUMENT (empty command), PERMISSION_DENIED
	rpc StartJob(StartJobRequest) returns (StartJobResponse);

	// StopJob terminates a job and all its child processes.
	// Errors: NOT_FOUND, PERMISSION_DENIED
	rpc StopJob(StopJobRequest) returns (StopJobResponse);

	// GetStatus returns the current state of a job.
	// Errors: NOT_FOUND, PERMISSION_DENIED
	rpc GetStatus(GetStatusRequest) returns (GetStatusResponse);

	// GetLogs streams the job's combined stdout/stderr from the beginning.
	// If follow=true, the stream continues until the job exits.
	// If follow=false, returns current output and closes.
	// Errors: NOT_FOUND, PERMISSION_DENIED
	rpc GetLogs(GetLogsRequest) returns (stream GetLogsResponse);
}
```

## Go API

```go
// Job defines a process to be executed with optional resource constraints.
type Job struct {
	Command        string          // Path to the executable
	Args           []string        // Command-line arguments
	ResourceLimits *ResourceLimits // Optional cgroup limits (nil = no limits)
}

// ResourceLimits specifies cgroups v2 resource constraints for a job.
type ResourceLimits struct {
	CPUMax         float64 // Fraction of CPU (0.5 = half a core, 2.0 = two cores)
	MemoryMaxBytes int64   // Memory limit in bytes
	IOMaxReadBPS   int64   // Max read throughput in bytes/sec
	IOMaxWriteBPS  int64   // Max write throughput in bytes/sec
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

// JobStatus represents the current state of a job.
type JobStatus struct {
	ID       uint64   // Unique job identifier
	State    JobState // Current state of the job
	ExitCode int      // Process exit code (only valid if State != JobStateRunning)
	Command  string   // Original command path
	Args     []string // Original arguments
}

// Worker manages the lifecycle of jobs and their associated resources.
type Worker struct {
	// private fields
}

// New creates a new Worker instance.
func New() *Worker

// Start executes a job and returns its unique ID.
func (*Worker) Start(Job) (id uint64, err error)

// Stop terminates a running job and all its child processes.
func (*Worker) Stop(id uint64) error

// Logs returns a reader for the job's combined stdout/stderr.
// If follow is true, the reader blocks waiting for new output until the job exits.
// The context controls cancellation; when cancelled, any blocked Read() returns immediately.
func (*Worker) Logs(ctx context.Context, id uint64, follow bool) (io.ReadCloser, error)

// Status returns the current state of a job.
func (*Worker) Status(id uint64) (JobStatus, error)
```

## Authentication & Authorization
The server and client authenticate using mTLS. Both sides present certificates signed by a shared CA. The server verifies client certificates and extracts the client identity (e.g., Common Name) from the client certificate for authorization.

Only TLS 1.3 is supported. This is to simplify cipher suite selection and it provides a strong set of suites to choose from. The server supports certificates in ECDSA with P-256, or RSA 2048-bit minimum. The CA, server cert, and client certs can be generated via a setup script.

Jobs are owned by the client identity (certificate Common Name) that created them. Only the owner of a job can perform subsequent operations on them. The gRPC server maintains a mapping of job IDs to owner identities and enforces ownership checks before delegating to the library.

## Implementation Details
### Process execution lifecycle
Each job is placed in a dedicated cgroup under `/sys/fs/cgroup/jobworker/job-{id}` and the resource limits are written to the appropriate files (`cpu.max`, `memory.max`, `io.max`). The mapping of users to jobs is maintained in memory by the gRPC server layer, keeping user identity out of the library.
The process is started directly in its cgroup using `clone3` with `CLONE_INTO_CGROUP` via `exec.Cmd`.
When stopping a job, the server will write to the `cgroup.kill` file for the job which will forcefully terminate all processes in the cgroup. Graceful shutdown (SIGTERM with a timeout) can be added in the future but is out of scope for the initial implementation.

Jobs are managed in a mutex-guarded map keyed by job ID. Job IDs are generated via atomic counter.
Each job tracks its `*exec.Cmd`, cgroup path, status (running/exited), and exit code.

On startup, the server kills any existing jobs and removes output directories to ensure a clean state. Jobs do not persist across server restarts.

### Output streaming
Process stdout and stderr are combined and written to a file at `/var/lib/jobworker/job-{id}/output`. This allows multiple concurrent clients to stream output independently, each reading from their own position in the file. Output is treated as opaque binary data with no encoding assumptions.

To avoid polling, the library uses a broadcast channel pattern. The writer holds a channel that is closed and replaced on each write, along with a version counter that increments on each write. Readers that reach EOF grab the current channel and version under a lock, then select on this channel and the context's Done channel. If the version has changed since the reader last checked, it skips waiting and retries immediately - this prevents a race where a signal is missed between reading EOF and starting to wait. When the process exits, the channel is closed one final time and a `done` flag is set.

**Client disconnection**: The `Logs()` method accepts a context. When a client disconnects, the gRPC layer cancels the context, which unblocks any waiting `Read()` call and returns `context.Canceled`. The reader closes its file handle and the goroutine exits cleanly. The underlying output file and job are unaffected.

**Job completion**: When the process exits, a `done` flag is set and the broadcast channel is closed. Readers wake up, see `done` is true, drain any remaining bytes from the file, then return EOF.

This approach keeps memory usage constant regardless of output size, supports late-joining clients reading from the start, and avoids busy-waiting.

### Resource control with cgroups
The library targets cgroups v2. Resource limits map as follows:
- `cpu_max`: Set `cpu.max` (e.g., 0.5 CPUs -> "50000 100000" for 50ms per 100ms period)
- `memory_max_bytes`: Set `memory.max`
- `io_max_read_bps` / `io_max_write_bps`: Set `io.max` with device major:minor numbers for all block devices

## Future work
- ListJobs: add functionality to list the jobs for a user
- DeleteJob: add functionality to remove a job from the system and cleanup its resources (output files, cgroup directories)
- Graceful shutdown: Send SIGTERM to the job process, wait for a termination window, and then send SIGKILL
- Per-user process isolation: Run jobs as distinct Linux users
- More sophisticated authorization scheme: role or attribute-based authorization can be more powerful and flexible
- Output retention policy: automatically clean up jobs after a period of time
- Resource defaults: server-side defaults for jobs that do not specify any
- Persistence: persist state to a durable store rather than the in-memory approach. Jobs should survive server restarts
- Observability instrumentation
- Per-device IO resource controls: Instead of one blanket setting for all devices, offer a way of configuring IO limits on a per-device basis
