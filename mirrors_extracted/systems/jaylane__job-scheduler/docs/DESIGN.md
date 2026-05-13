---
authors: Jason Lane (jason.j.lane@gmail.com)
state: draft
---

# What

A prototype job worker service that provides an API to run arbitrary Linux processes.

# Why

To level set my engineering skill, and spend 2 weeks working on a project to get a feel for the day to day of an engineer at Teleport.

# Details

This project will consist of 3 main parts per the level 4 requirements spec laid out in the [challenge document](https://github.com/gravitational/careers/blob/main/challenges/systems/challenge-1.md).

## Library

The library will consist of 3 Go packages that will be used by the API and the CLI to run the processes.

### Worker package

The worker will be the top level struct in the library. Workers will be able to injest a configuration object which will define:
* Log output directory
* Log file name
* Log chunk size
* Logger

Workers will have 2 properties

- a logger
- a map of jobs (a map of jobID (uuid) to a pointer to the specific job.)

The worker struct will have 4 methods:

- StartJob - (command string, args []string) (jobID string, err error) 
    - Starts a job 
    - command is a program executable available via PATH or an absolute path to the program 
    - args are arguments to supply to the command
    - jobID is a string version of a uuidv4 of the created job
- StopJob (jobID string) (err error)
    - jobID is the ID of the job to stop
- GetJobStatus (jobID string) (status job.Status, err error)
    - jobID is the ID of the job to get status on
    - status is the status of the job 
- StreamJobOutput(ctx context.Context, jobID string) (outchan chan string, err error)
    - ctx context
    - jobID is the ID of the job to stream the output of
    - chan for streaming output from the log

    The worker will write the output of the job's process (stderr/stdout) on the server via a log file associated to the job by its ID (ex: 975b2d14-e567-4e22-92e4-eebefe6d8ed7.log). For streaming the logfile to the client the worker's logger will incorporate a filewatcher that will listen for file events on the specific job's logfile via inotifywait. As the logfile remains open and continues to be modified the logger will stream the incoming content to the client via a channel.

#### Trade-off
For the interest of time each worker will store the logfile locally on the server. The log files will persist through the life of the server. This would not be ideal in a production instance. In a real world implemention I would store logging data in an s3 bucket or via a vendor logging solution such as Datadog.

### Job package

The job package will define the structs containing the information about the specific processes spawned by the worker.

Example of the structs that will make up the Job package

```golang
type Process struct {
	PID int `json:"pid"`
	ExitCode int `json:"exit_code"`
	Status Status `json:"status"`
}

type Status int

const (
	UNKNOWN Status = iota
	RUNNING
	STOPPED
    TERMINATED
    COMPLETED
)

func (s Status) String() string {
	switch s {
	case RUNNING:
		return "running"
	case STOPPED:
		return "stopped"
    case TERMINATED
        return "terminated"
    case COMPLETED
        return "completed"
	default:
		return "unknown"
	}
}

func (s Status) EnumIndex() int {
	return int(s)
}

type Job struct {
    ID string 
    Cmd *exec.Cmd 
    Process *Process
    StartTime time.Time
    EndTime time.Time
}
```

### cgroups package

cgroups implementation will be added via a cgroup package that will be used to create the necessary directory structure in the server's filesystem. cgroups will be separated by jobID upon starting of the job and any child procs will be added by piping the pids the cgroups procs config. Resource limits for CPU, Memory, and BlockIO will be hardcoded for this project: 

```golang
const (
    CPUPeriodUs = 100000
    CPULimit = 2.0
    MemoryLimit = "1G"
    BlkioMajorMinor = "8:16"
    BlkioWriteLimit = BlkioMajorMinor + " 20971520"
    BlkioReadLimit = BlkioMajorMinor + " 41943040"
)
```
- memory:
    ```bash
        echo 1G > /sys/fs/cgroup/memory/jobs/{jobID}/memory.limit_in_bytes
    ```
- cpu:  
    ```bash 
       echo 100000 > /sys/fs/cgroup/cpu/jobs/{jobID}/cpu.cfs_period_us
       echo 200000 > /sys/fs/cgroup/cpu/jobs/{jobID}/cpu.cfs_quota_us
    ```
- disk_io:
    ```bash 
        echo 8:16 41943040 > /sys/fs/cgroup/blkio/jobs/{jobID}/blkio.throttle.read_bps_device
        echo 8:16 20971520 > /sys/fs/cgroup/blkio/jobs/{jobID}/blkio.throttle.write_bps_device
    ```

#### Trade-off 
In a production environment the resource limits could be configurable but for this project I will hardcode limits to reduce scope.


## API
The API will be a gRPC server that will be responsible for authentication, authorization and interacting with the library to execute Worker methods:

https://github.com/jaylane/job-scheduler/blob/acf2ea9742674d9d15b7dda6df87cc97975cb756/proto/worker.proto#L1-L68


## CLI
The CLI will utilize [cobra](https://github.com/spf13/cobra) for ease of development. The CLI will have a gRPC client that will be able to interact with the API to start/stop/get status/stream output of jobs in their local shell.


example: 
```sh
jsched-cli start "/bin/ls" "-l"

JobID: aeba5ba9-e95a-455b-b97b-31a2d98c45ab started

jsched-cli stop "aeba5ba9-e95a-455b-b97b-31a2d98c45ab"

jsched-cli status aeba5ba9-e95a-455b-b97b-31a2d98c45ab

JobID: aeba5ba9-e95a-455b-b97b-31a2d98c45ab running/stopped/terminated

jsched-cli stream aeba5ba9-e95a-455b-b97b-31a2d98c45ab

2024/08/11 22:17:25 Running command /bin/ls -l
2024/08/11 22:17:25 total 16
2024/08/11 22:17:25 -rw-r--r--@ 1 jasonlane  staff  1065 Aug  7 14:04 LICENSE
2024/08/11 22:17:25 drwxr-xr-x  3 jasonlane  staff    96 Aug  7 14:04 docs
2024/08/11 22:17:25 -rw-r--r--  1 jasonlane  staff   573 Aug 11 22:17 main.go
2024/08/11 22:17:25 drwxr-xr-x  3 jasonlane  staff    96 Aug  7 14:06 proto
```


## Transport

TLS 1.3 has been chosen for secure transport of communication between client and server for a few important reasons:

- Improved Security: TLS 1.3 removes outdated cryptographic algorithms, making connections more secure against modern threats.

- Performance Gains: The handshake process is simplified and faster (one round trip), reducing latency and improving performance, especially in high-traffic environments.

- Forward Secrecy: TLS 1.3 enforces forward secrecy, ensuring that even if a key is compromised, past communications remain secure.

- Simplified Configuration: With fewer options and mandatory security features, it reduces the chances of misconfigurations.

### Cipher Suites

For this project I will support the following tls 1.3 cipher suites, opting to be more secure rather than have wider compatibility.

```golang
TLS_AES_128_GCM_SHA256       uint16 = 0x1301
TLS_AES_256_GCM_SHA384       uint16 = 0x1302
TLS_CHACHA20_POLY1305_SHA256 uint16 = 0x1303
```

#### Trade-off 
Some legacy systems or older network devices would not be compatible with these newer cipher suites, in a real world application this would have to be discussed to find a happy medium between security and compatibility.

## mTLS

Authentication and Authorization will be done via mTLS. I will create a bash script to generate the following certificates and store them in the repository for this project. 

* Client CA private key and signed cert
* Server CA private key and signed cert
* Client private key and CSR
* Server private key and CSR
* Server Signed Cert via CSR & CA noted above
* Client Signed Cert via CSR & CA noted above

The API Server will load its key and certificate as well as the client's CA via a configuration object in its main.go file (`/cmd/api/main.go`). 

```golang
type ServerConfig struct {
	Address string
	Certificate string
	Key         string
	ClientCA          string
}
```

The client will follow the same pattern instead loading the server's CA in its main.go file. (`/cmd/client/main.go`).

```golang
type ClientConfig struct {
	ServerAddress string
	Certificate string
	Key         string
	ServerCA          string
}
```

### Certificates
* X.509
* Signature Algorithm: EdDSA with ED25519 scheme
* Public Key Algorithm: EdDSA with ED25519 scheme
* ED25519 Public-Key: (256 bit)
* roleOid 1.2.840.10070.8.1 = ASN1:UTF8String 

After doing some more research I decided to move away from RSA which was a kneejerk reaction because its what I've used in the past. 
The reasons I decided to switch to EdDSA:

- Offers excellent security and performance, even on less powerful hardware.
- Simplicity: Uses fixed-size keys (256 bits) and avoids the pitfalls of ECDSA’s curve parameters.

#### Trade-off 
The downside to EdDSA is that is does not have universal adoption, but like the trade-off above with the cipher suites for this project I'm opting for performance and security over widespread compatibility.

There will be 2 roles as far as authorization is concerned admin & user. The client certificates will have these as X.509 v3 extensions. The API will use middleware to either authorize or reject the request based on the incoming certificate's role. Following the example from the spinnaker [docs](https://spinnaker.io/docs/setup/other_config/security/authentication/x509/#creating-an-x509-client-certificate-with-user-role-information). The roles must be informed when the CA signs the CSR so the extension attribute roleOid 1.2.840.10070.8.1 = ASN1:UTF8String must be requested in the CSR when the client cert is created. The data following UTF8String: is encoded inside of the x509 cert under the given OID. 

Usernames will be added to the X.509 certificates by using the Subject field under the CommonName RDN.

admin role - authorized to interact with all rpcs

client-admin-ext.conf
```conf
subjectAltName=DNS:localhost
  [ v3_req ]
  1.2.840.10070.8.1 = ASN1:UTF8String:admin
```

user role - authorized to interact with all but StopJob rpc

client-user-ext.conf
```conf
subjectAltName=DNS:localhost
  [ v3_req ]
  1.2.840.10070.8.1 = ASN1:UTF8String:user
```

#### Trade-off 
In a real world scenario these certificates could be managed by secrets or using something like Vault.

