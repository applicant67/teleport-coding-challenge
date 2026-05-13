# Job Worker Design Doc

Design doc for a prototype job worker service that provides an API to run arbitrary Linux processes.

## Level 2 requirements

- Library: Worker library with methods to start/stop/query status and get the output of a job.
- API:
    - HTTPS API to start/stop/get status of a running process.
    - Use mTLS authentication and verify client certificate. Set up strong set of cipher suites for TLS and good crypto
      setup for certificates. Do not use any other authentication protocols on top of mTLS.
    - Use a simple authorization scheme.
- CLI:
    - CLI should be able to connect to worker service and schedule several jobs.
    - CLI should be able to query result of the job execution and fetch its output.

## Job object

- `ID string` unique identifier of a job
- `Command string` Linux command to be executed
- `Arguments []string` list of arguments for the command
- `Status string` current status of the job
- `OutputBuffer bytes.Buffer` links stderr and stdout of the job process, it will be hidden from API response
- `Output string` output of the job process, converted from `OutputBuffer` when status is queried
- `Error error` startup error, termination error, or exit error
- `ExitCode integer` exit code of the job process

## Job status

- `Running` Job process started
- `Exited` Job finished successfully or exited with error or force stopped

## Job handling

- Start the process using `exec.Cmd` in a separate go-routine
- Utilize go-channels to communicate with the go-routine
- The go-routine should exit when the job finishes or force stop requested

## Library

- `CreateJob(command)`
    - Generate a unique ID
    - Initialize a job object using the ID
    - Save the reference to the job into memory for future access
    - Start the job in a separate go-routine
    - Return the job object when it successfully starts
    - Return error if job fails to start
- `StopJob(jobID)`
    - Fetch the job object from memory by job ID
    - Get the current status of the job and check if it is running or already stopped
    - Stop the job process if it is running using `Cmd.Process.Kill()`. It returns error if termination fails and keep status `Running`. Otherwise, send a confirmation on success and change status to `Exited`.
- `GetJobStatus(jobID)`
    - Fetch the job object from memory by job ID
    - Convert current output of the job from buffer to string
    - Return the details of a job along with output and error

## API

- `POST /create`

Request:
```json
{
  "command": "bash",
  "args": ["-c", "pwd"]
}
```

Response:
```json
{
  "id": "c2l1qbucie6hpdufnung",
  "command": "bash",
  "args": ["-c", "pwd"],
  "status": "CREATED",
  "output": "",
  "error": "",
  "exitCode": 0
}
```

- `GET /status/{id}`

Response:
```json
{
  "id": "c2l1qbucie6hpdufnung",
  "command": "bash",
  "args": ["-c", "pwd"],
  "status": "EXITED",
  "output": "/Users/das/Downloads/teleport/job-worker\n",
  "error": "",
  "exitCode": 0
}
```

- `PUT /stop/{id}`

Response:
```json
{
  "message": "Job c2l1qbucie6hpdufnung terminated"
}
```

## CLI

- `job-worker [server/client]` combined CLI for running server and client
- `job-worker server [--port] [--cacert] [--cert] [--key]` run the API server
- `job-worker client [create/stop/status] [--address] [--cacert] [--cert] [--key]` run client to communicate with the API server
- `job-worker client create [--command]` request a new job
- `job-worker client stop [--id]` stop a job
- `job-worker client status [--id]` fetch the status of a job

## mTLS

- Use `openssl` to generate self-signed certificates
- Use `RSA 2048` key and `sha256` digest
- Shell script to generate certificates
- Steps:
    - Generate a root CA certificate (`ca.crt`, `ca.key`)
    - Generate server certificate signed by the CA (`server.crt`, `server.key`)
    - Generate client certificate signed by the CA (`client.crt`, `client.key`)

## Authorization

- Simplified RBAC style authorization
- Permissions:
    - Read and write: Required for creating and stopping a job
    - Read only: Required for querying status of a job
- Two hard coded users to mock database
    - Username `alice`: has read and write permissions
    - Username `bob`: has read only permission
- Hard coded username to permission mapping
- Pass username through client certificate using `Common Name` i.e. each user will have separate client certificate

## Tests

- Test the library
- Test RBAC authorization
- Test TLS connection

## Third party libraries

- [spf13/cobra](https://github.com/spf13/cobra) for CLI
- [gorilla/mux](https://github.com/gorilla/mux) for routing
- [rs/xid](github.com/rs/xid) for generating globally unique job ID
- [sirupsen/logrus](github.com/sirupsen/logrus) for logging

## TODOs

- Terminate all jobs before shutting down the server
- Support graceful termination of jobs along with force stop
- Retry job termination if failed
- Use DB to store users
- Save job details in DB instead of memory
- Use job queue to schedule the job based on resource availability instead of directly starting the job 
- Job chaining (workflow) using dependency graph
- Check for vulnerabilities before running a job
