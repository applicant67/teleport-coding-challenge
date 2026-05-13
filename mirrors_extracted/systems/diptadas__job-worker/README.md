# Job Worker

A prototype job worker service that provides an API to run arbitrary Linux processes.

## Install

```shell
$ git clone https://github.com/diptadas/job-worker.git
$ cd job-worker
$ goimports -l -w ./
$ go install
$ which job-worker
/Users/das/go/bin/job-worker
```

## Run Tests

```shell
$ go test -v ./...
```

## Library usage

- Create a job

```go
request := CreateJobRequest{
    Command: "echo",
    Args:    []string{"hello world"},
}

job, err := CreateJob(request)
```

- Stop a job

```go
job, err := StopJob("c2vvjbucie6ju0ou7j5g")
```

- Get job status

```go
job, err := GetJobStatus("c2vvjbucie6ju0ou7j5g")
```

## CLI

```shell
$ job-worker -h
CLI for the job-worker

Usage:
  job-worker [command]

Available Commands:
  client      Client for job-worker service
  help        Help about any command
  server      Start the job-worker API server

Flags:
  -h, --help   help for job-worker

Use "job-worker [command] --help" for more information about a command.
```

```shell
$ job-worker server -h
Start the job-worker API server

Usage:
  job-worker server [flags]

Flags:
      --ca string     CA cert file (default "ssl/ca.crt")
      --cert string   Server cert file (default "ssl/server.crt")
  -h, --help          help for server
      --key string    Server key file (default "ssl/server.key")
  -p, --port int      Port of API server (default 8443)
```

```shell
$ job-worker client -h
Client for job-worker service

Usage:
  job-worker client [command]

Available Commands:
  create      Create a new job
  status      Get status of a job
  stop        Stop a job

Flags:
      --address string   Address of the API server (default "https://localhost:8443")
      --ca string        CA cert file (default "ssl/ca.crt")
      --cert string      Client cert file
  -h, --help             help for client
      --key string       Client key file

Use "job-worker client [command] --help" for more information about a command.
```

## API endpoints

- `POST /create` Creates a new job and returns the job ID

```json
{
  "command": "bash",
  "args": ["-c", "pwd"]
}
```

- `GET /status/{id}` Returns the status of a job
- `PUT /stop/{id}` Stops a job process

## Generate TLS certificates

TLS certs are already generated in the `ssl` directory. New certs can be generated using `gen-cert.sh` script.

## Run API server

### Run in host

```shell
$ job-worker server
INFO[0000] Server started at :8443
```

### Run in Docker (recommended)

```shell
$ docker build -t diptadas/job-worker .
$ docker run -p 8443:8443 -it diptadas/job-worker
```

## Send Requests

### Using CLI

```shell
$ job-worker client create --cert ssl/client-alice.crt --key ssl/client-alice.key --cmd pwd | jq .
{
  "id": "c353t906n88jjt7idsmg",
  "status": "RUNNING",
  "output": "",
  "request": {
    "command": "pwd",
    "args": null
  },
  "error": null
}

$ job-worker client status --cert ssl/client-alice.crt --key ssl/client-alice.key --id c353t906n88jjt7idsmg | jq .
{
  "id": "c353t906n88jjt7idsmg",
  "status": "EXITED",
  "output": "/Users/das/Downloads/job/teleport/job-worker\n",
  "request": {
    "command": "pwd",
    "args": null
  },
  "error": null
}
```

### Using cURL

```shell
$ curl -X GET https://localhost:8443/status/c353t906n88jjt7idsmg \
  --cert ssl/client-alice.crt \
  --key ssl/client-alice.key \
  --cacert ssl/ca.crt | jq .
{
  "id": "c353t906n88jjt7idsmg",
  "status": "EXITED",
  "output": "/Users/das/Downloads/job/teleport/job-worker\n",
  "request": {
    "command": "pwd",
    "args": null
  },
  "error": null
}
```
