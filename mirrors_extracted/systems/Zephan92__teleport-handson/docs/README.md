# Teleport Hands-On: Job Worker Service

A prototype job worker service providing a gRPC API for remote execution of isolated Linux processes with real-time log streaming.

This project implements a secure, gRPC-based Job Worker Service capable of executing arbitrary Linux processes with resource isolation (Cgroups) and real-time output streaming. It includes a central server and a CLI client, communicating via mTLS.

## Features Implemented

- **Core Worker Library**: Manages Linux process lifecycle (Start, Stop, Wait).
- **Resource Isolation**: Enforces CPU, Memory, and I/O limits using Cgroups v2.
- **Process Orchestration**: Uses `ptrace` to safely jail processes before execution.
- **Log Streaming**: Implements a thread-safe, pull-based iterator pattern for efficient concurrent log access.
- **gRPC Server**: Supports `Start` and `Stop` operations with panic recovery and graceful shutdown.
- **Integration Testing**: Verifies lifecycle, error handling, and kernel-level resource enforcement (OOM killing).

## Documentation

- **[Challenge Requirements](challenge-1.md)**: The original problem statement and requirements.
- **[Design Document (RFD)](0001-job-worker-service.md)**: The technical design, architecture, and API specification.

## Project Structure

- **`api/v1`**: Protobuf definitions and generated Go code.
- **`cmd/`**: Entry points for the Server and CLI binaries.
- **`internal/worker`**: Core library logic for process management and Cgroups.
- **`internal/server`**: gRPC server implementation.
- **`internal/cli`**: CLI command logic and configuration.
- **`scripts/`**: Helper scripts for environment setup.
- **`docs/`**: Design documents and specifications.

## Prerequisites

- **Go**: 1.22+
- **Linux Kernel**: 5.14+ (Required for `cgroup.kill` support)
- **Protoc**: Protocol Buffers compiler (for regenerating code)
- **Make**: For running build commands

## Build Instructions

To build the binaries for your local OS:
```bash
make build
```

To force a Linux build (required for the worker server functionality if building on non-Linux):
```bash
make build-linux
```

To build the CLI for Windows (if developing on Windows):
```bash
make build-windows
```

## Usage

### Environment Setup
Before running the server, ensure your Linux environment is configured for Cgroup v2 delegation. This script checks the kernel version and enables controller delegation.

```bash
sudo ./scripts/setup-cgroups.sh
```

### Start the Server
The server requires root privileges to manage cgroups.

```bash
sudo ./bin/worker-server
```

## CLI Usage
The `worker-cli` allows you to interact with the server. It uses mTLS for authentication.

### Global Flags:

- `--addr`: Server address (default "localhost:8080"). Env: `WORKER_ADDR`
- `--ca`: Path to CA certificate (default "certs/ca.crt"). Env: `WORKER_CA`
- `--cert`: Path to client certificate (default "certs/alice.crt"). Env: `WORKER_CERT`
- `--key`: Path to client key (default "certs/alice.key"). Env: `WORKER_KEY`
- `--json`: Output response in JSON format.
- `--quiet` (`-q`): Suppress informational output (useful for scripting).

### Commands:
#### Start a Job
Starts a new job. Use `--` to separate flags from the command to run.

##### Flags:
- `--follow` (`-f`): Stream logs after starting the job.
```bash
# Start a job
./bin/worker-cli start -- echo "Hello World"

# Start and follow logs immediately
./bin/worker-cli start -f -- sh -c "while true; do echo ping; sleep 1; done"

# Specify working directory
./bin/worker-cli start --dir /tmp -- ls -la
```

#### Stop a Job
```bash
# Terminates a running job.
./bin/worker-cli stop <job_id>
```

#### Status Logs
```bash
# Retrieves status, exit code, and resource usage.
./bin/worker-cli status <job_id>
```

#### Stream Logs
```bash
# Streams stdout and stderr from a running or completed job.
./bin/worker-cli stream <job_id>
```

#### Version
```bash
# Prints the CLI version.
./bin/worker-cli version
```

## Development

Generate Protobuf code:
```bash
make proto
```

Generate Certs for Int Tests:
```bash
make certs
```

### Testing Strategy & Cgroup Isolation

The project employs a dual-mode testing strategy to ensure robustness while maintaining developer velocity:

1.  **User-Mode Integration Tests (`make test`)**:
    *   These tests run without root privileges.
    *   They use `t.TempDir()` to simulate the Cgroup filesystem structure.
    *   **Note**: Since standard filesystems (ext4/tmpfs) behave differently than Cgroupfs (specifically regarding `rmdir` on non-empty directories), the worker library includes a safety check. It detects if the target is *not* a real Cgroup filesystem and safely falls back to recursive deletion. This ensures logic verification without requiring a VM or root access for every run.

2.  **Root-Mode Integration Tests (`make test-root`)**:
    *   These tests require `sudo` and interact with the real Linux Cgroup v2 hierarchy (`/sys/fs/cgroup`).
    *   They verify kernel-level enforcement of resource limits (e.g., OOM killing) and atomic process tree termination.
    *   In this mode, the test cleanup logic strictly enforces `unix.Rmdir`, matching production behavior.
 

**3. Full Coverage**
To run **all** tests and generate a comprehensive coverage report:
```bash
make test-complete
```