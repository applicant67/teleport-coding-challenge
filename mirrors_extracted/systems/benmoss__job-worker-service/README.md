# Job worker service

## Getting started

This project relies on Linux cgroups v2, and as such will only work on Linux. There are lots of tools for running a Linux VM if you're developing on another OS. I use [lima](https://lima-vm.io/), but setup instructions for that are beyond the scope of this README.

### Dependencies
- [pixi](https://pixi.prefix.dev)

## Building

Build the server binary:
```bash
pixi run build-workerd
```

## Running the server

The server requires `sudo` to create and manage cgroups:

```bash
sudo ./bin/jobworkerd \
  -addr :8080 \
  -cgroup-root /sys/fs/cgroup/jobworker \
  -log-dir /var/lib/jobworker
```

### Configuration flags
- `-addr`: gRPC server address (default: `:8080`)
- `-cgroup-root`: Root directory for job cgroups (default: `/sys/fs/cgroup/jobworker`)
- `-log-dir`: Directory for job output logs (default: `/var/lib/jobworker`)

### Example client usage

Using `buf curl`

```bash
# Start a job
pixi run buf curl --protocol grpc --http2-prior-knowledge \
  -d '{"command": "/bin/sleep", "args": ["10"]}' \
  http://localhost:8080/jobworker.v1.JobWorkerService/StartJob

# Get job status
pixi run buf curl --protocol grpc --http2-prior-knowledge \
  -d '{"job_id": 1}' \
  http://localhost:8080/jobworker.v1.JobWorkerService/GetStatus

# Stream job logs
pixi run buf curl --protocol grpc --http2-prior-knowledge \
  -d '{"job_id": 1, "follow": true}' \
  http://localhost:8080/jobworker.v1.JobWorkerService/GetLogs

# Stop a job
pixi run buf curl --protocol grpc --http2-prior-knowledge \
  -d '{"job_id": 1}' \
  http://localhost:8080/jobworker.v1.JobWorkerService/StopJob
```

## Regenerating protobuf code

If you modify [proto/jobworker/v1/service.proto](proto/jobworker/v1/service.proto), regenerate the Go code:

```bash
pixi run buf generate
```

## Running tests
You can run `sudo go test ./...` to run the Go tests. `sudo` is required since by default a non-root user does not have permission to create cgroups.
