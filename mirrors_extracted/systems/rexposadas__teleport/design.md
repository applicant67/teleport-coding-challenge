---
authors: Rex Posadas (rexposadas@gmail.com)
---

# What we are building

Take-home challenge for Level 4 Systems Engineer.

This documentation is a design for a job worker service and its associated CLI.

The implementation is meant to be minimal while meeting all challenge requirements. 
Tradeoffs are noted in the `tradeoffs` section of this document.

# Parts of the challenge:

## Worker Library

1. Has methods to start/stop/query status and stream logs from the jobs.
1. Resource controls of CPU, memory and disk IO using cgroups.

## API

1. gRPC API to start/stop/get status/stream output of a running process.
1. Use mTLS for authentication.
1. Has a simple authentication scheme.

## Client

1. CLI connects to worker service to start, stop, get status, and stream output of a job.
2. Multiple clients can connect to the service.


# Design

## Worker Library and Server

The server will use a library to: 

1. Start a job. Starting a job returns a UUID which is the job id. This ID can be used to do further
2. Stop a job. I place the process in the cgroup and kill the cgroup itself. That will remove all processes spawned by a job.
3. Query job.  Returns the status of a job. 
4. Provide streaming of job logs.

Service should be able to handle multiple clients. 

Sends SIGKILL to the process. SIGKILL in

### Tradeoffs and Production 
- Initially using SIGTERM, then SIGKILL as a last resort is probably the better way to do things. But for this 
  challenge, SIGKILL is sufficient.

### cgroups

`cgroups v2` will be used for resource control. 

- Each job will run in its own cgroup.
- Resource usage will be isolated per job.
- There will be job_id grouping. for example: 
  - /sys/fs/cgroup/teleport/<job_id> -- the job_id is the UUID provided by the server when the job was created.
- Use cpu.max, memory.max and io.max controllers to set resource values.

The hardcoded resource limits are below. For this challenge, I will assume the jobs are relatively light: 
- CPU: 0.5 
- Memory: 256MB
- IO: 1MB/s

For simplicity and for this challenge, I will assume there is only one block device /dev/sda with major/minor 8:0. 

#### Placing jobs in cgroups

- os.Open opens /sys/fs/cgroup/teleport/<job_id>, resulting in a file descriptor.
- When running exec.Cmd, use UseCgroupFD and CgroupFD to specify that the process should be placed in a cgroup.


#### Trade offs and the Production environment
- It is preferable to dynamically configure resource limits based on job and system state. A first step in doing 
  this is providing a way to adjust resource limits when jobs are created.
- In production, we probably want orchestrators, like systemd, which manages cgroups. 

## API

### Using gRPC for communication. 

- Start: Starts a job with hardcoded resource limits.
- Stop: Stops a job.
- GetStatus: Fetches status of a job
- GetLogs: Stream logs of a running job.

```
service TeleportService {
    rpc Start (StartRequest) returns (StartResponse);
    rpc Stop (StopRequest) returns (StopResponse);
    rpc GetStatus (GetStatusRequest) returns (GetStatusResponse);
    rpc GetLogs (GetLogsRequest) returns (stream GetLogsResponse);
}

message StartRequest {
    string command = 1;
    repeated string args = 2;
}

message StartResponse {
    string job_id = 1; // a job UUID
}

message StopRequest {
    string job_id = 1;
}

message StopResponse {
    // ideally, add a boolean to see if stopping was successful
    // maybe even a message. But not necessary for the challenge.
}

message GetStatusRequest {
    string job_id = 1;
}

message GetStatusResponse {
    Status status = 1;
}

enum Status{
    STATUS_UNSPECIFIED = 0; 
    STATUS_RUNNING = 1;
    STATUS_STOPPED = 2; // User stopped the job
    STATUS_EXITED = 3;  // Job exited
    STATUS_UNKNOWN = 4; // For anything we don't recognize
}

message LogsRequest {
  string job_id = 1;
}

message LogsResponse {
  bytes data = 1; // log data (both stdout and stderr)
}
```

### TLS

API will use TLS 1.3, which is the latest version. 1.3 offers improved security and performance gains the previous 
versions. 

Recommended curves: 
- X25519 : Performant and secure.
- P-256 : Widely supported and TLS 1.2 compatible.
- P-384 : Higher level of security. Good for sensitive data.
- P-521 : Very high security.

Recommended Cipher suites. These are performant and offers good security:
- TLS_AES_128_GCM_SHA256
- TLS_CHACHA20_POLY1305_SHA256
- TLS_AES_256_GCM_SHA384

### Authentication

mTLS will be used for authentication. Both server and client provides certificates in order to 
communicate. For testing, premade certificates will be provided.

Users are authorized to use the service by its certification being validated. Validation is done by looking at the 
Common Name (CN), which is part of the subject.

The ownership model would be such that there is a mapping between the user in the CN and job id.  This would 
indicate what jobs a user can manage.

The CLI will look for certificates in the `~/.certs`.  These are the necessary files:

Certificates will use ECDSA as the signature algorithm.

Server:
- `~/certs/ca.pem`
- `~/certs/service.pem`
- `~/certs/service-key.pem`

Client:
- `~/certs/ca.pem` 
- `~/certs/client-<id>.pem` 
- `~/certs/client-key-<id>.pem`

To create certificates, I would use `make` (Makefile) targets to generate the certificates.
For examples `make cert name=adam` will create the certificates for user adam.


#### Tradeoffs
- The root CA are self-signed. Production applications should use a well-trusted CA. 
- There is a single Certificate Authority (CA) for both the client and server.  Having different CA for server 
  and client is more secure. But for this challenge, a single CA is sufficient.
- The CLI only looks at `~/certs` for certificates. That's hardcoded.  Adding that as a CLI option would be easy enough, but didn't think it was necessary for this challenge. 



### Streaming logs

We will use a pub/sub model to stream logs to multiple clients.  Logs, stored in-memory, are
sent to the subscribers.

- Server stores logs in an in-memory buffer (publisher).
- Clients wanting logs subscribes to receive messages. (subscriber)
- Common to pub/sub models, when a subscriber starts receiving logs, it gets what's in the buffer followed by the 
  feed.
- Both stdout and stderr will be logged.
- Go primitives used:
  - channels
    - a channel is created for every client.  A client will listen to the channel for logs.
  - go routines
    - a go routine collects logs from stderr and stdout and puts it in an in-memory buffer. 
    - A go routine is created to read logs from in-memory buffer sending them to the proper channels so that 
      clients gets the logs. 
  - a mutex is used in the in-memory buffer so that we can sync writing and reading from it.

That's a general overview. Something might come up during the coding process that would change this plan. 

## Client

The client/cli manages jobs and certificates.

The CLI will be implemented using the [Cobra library](https://github.com/spf13/cobra) . 

### Start the server

> $ teleportd


###  Managing jobs

Below is an example of how to manage a job.

>`teleport start` 

>`teleport status <job_id>`  

>`teleport logs <job_id>` 
 
> `teleport stop <job_id>`


1. `start` returns a `job id`, which is a UUID
2. `status` fetches current status and resource usage.
3. `logs` displays the logs of the job
4. `stop` kills a job


# Tradeoffs

1. Logs are stored in-memory. In production, this is not recommended because of memory limits. Writing to files 
   is better. 
1. Geting a list of running jobs will not be implemented.
1. Setting resource controls such as CPU and memory usage per job will be hardcoded. Ideally, it can be set via the 
   CLI. Since we are using cgroups, we should be able to set different resource allocations per job.
1. Hardcoded cgroup settings has downsides: 
   1. This is not an efficient way to work with available resources. One can have jobs that 
         will not make use of the allocated resources. While, other jobs might require more. 
   1. If a job exceeds resources allocated unexpected behaviour might occur. This issue _might_ be hard to debug.
1. There is no scheduler.
1. The system will not scale and is not highly performant. 
1. This system is not highly available.




# Testing

1. Authentication: User actions with valid and invalid credentials. 
2. Server can handle multiple clients.
3. Start/stop/get status of jobs




