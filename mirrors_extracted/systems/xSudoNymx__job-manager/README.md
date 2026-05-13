# Job Manager Service

A prototype job manager service that provides a gRPC API to run arbitrary Linux processes.

## Project Structure

```
job-manager/
├── cmd/
│   ├── jobd/            # gRPC server entry point
│   └── jobctl/          # CLI client entry point
├── pkg/                 
│   ├── job/             # Job lifecycle management
│   ├── manager/         # Job coordination
│   ├── output/          # Output capture and streaming
│   └── client/          # gRPC client wrapper
├── internal/            
│   ├── server/          # gRPC server and handlers
│   └── auth/            # mTLS and authorization
├── proto/job/v1         # Protocol buffer definitions
└── certs/               # TLS certificates
```

## Quick Start

### Prerequisites

- Go 1.25+

### Build

From the root of the project, run:

```bash
go build -o jobd ./cmd/jobd
go build -o jobctl ./cmd/jobctl
```

This will produce the server (`jobd`) and client (`jobctl`) executables.

### Generate Certificates

Certificates are provided for testing authentication. By default, they use `localhost` and `127.0.0.1` for DNS and IP
SANs. For remote testing, regenerate with your DNS and IP:

```bash
Usage: ./certgen.sh [certs_dir] [dns_name] [ip_addr]

# Default: localhost, 127.0.0.1
./certgen.sh

# Custom DNS name
./certgen.sh certs myserver.local

# Custom DNS and IP
./certgen.sh certs myserver.local 192.168.1.10
```

Generated certificates:

| File                       | Purpose                    |
|----------------------------|----------------------------|
| `ca.crt`, `ca.key`         | Certificate Authority      |
| `server.crt`, `server.key` | Server certificate         |
| `admin.crt`, `admin.key`   | Admin client (OU=Admin)    |
| `read.crt`, `read.key`     | Read-only client (OU=Read) |

### Run the Server

```bash
./jobd
```

#### Server Flags

| Flag     | Default            | Description                 |
|----------|--------------------|-----------------------------|
| `--addr` | `:9000`            | gRPC listen address         |
| `--cert` | `certs/server.crt` | Server certificate          |
| `--key`  | `certs/server.key` | Server private key          |
| `--ca`   | `certs/ca.crt`     | CA certificate              |
| `--data` | (temp dir)         | Data directory for job logs |

### Use the CLI

```bash
# Start a job
./jobctl start -- /bin/echo "Hello, World!"

# Check status
./jobctl status <job-id>

# Stream output
./jobctl stream <job-id>

# Stop a job
./jobctl stop <job-id>
```

## Configuration

The CLI can be configured via flags, environment variables, or config file.

**Priority:** flags > environment variables > config file > defaults

### Flags

```bash
./jobctl --addr localhost:9000 --cert certs/admin.crt --key certs/admin.key --ca certs/ca.crt start -- /bin/echo hello
```

### Environment Variables

```bash
export JOBCTL_ADDR=localhost:9000
export JOBCTL_CERT=certs/admin.crt
export JOBCTL_KEY=certs/admin.key
export JOBCTL_CA=certs/ca.crt

./jobctl start -- /bin/echo hello
```

### Config File

Use `jobctl.yaml` in the current directory or `~/jobctl.yaml`:

```yaml
addr: localhost:9000
cert: certs/admin.crt
key: certs/admin.key
ca: certs/ca.crt
# server-name: localhost  # Override for TLS verification
# timeout: 30s
```

## Security

### TLS Configuration

- TLS 1.3 only
- mTLS with client certificate verification
- Certificates signed by a private CA

### Authorization

Roles are derived from the client certificate's Organizational Unit (OU):

| Role  | Permissions                 |
|-------|-----------------------------|
| Admin | start, stop, status, stream |
| Read  | status, stream              |

## Design

See [DESIGN.md](DESIGN.md) for detailed design documentation.