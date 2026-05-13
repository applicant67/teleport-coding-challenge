# Certificate Management

TLS certificates for the Job Worker Service using [cfssl](https://github.com/cloudflare/cfssl).

## Quick Start

```bash
# Install cfssl (macOS)
brew install cfssl

# Generate all certificates
make all

# Verify certificates
make verify-chain
```

## Certificate Hierarchy

```
Job Worker CA (ECDSA P-256, 10 years)
├── server.pem      (CN=worker-server, OU=server, 1 year)
├── client1.pem     (CN=client1, OU=user, 1 year)
├── client2.pem     (CN=client2, OU=user, 1 year)
└── admin.pem       (CN=admin, OU=admin, 1 year)
```

## Generated Files

| File | Description |
|------|-------------|
| `ca.pem` / `ca-key.pem` | CA certificate and private key |
| `server.pem` / `server-key.pem` | Server certificate (SANs: localhost, 127.0.0.1) |
| `client1.pem` / `client1-key.pem` | Regular user certificate |
| `admin.pem` / `admin-key.pem` | Admin certificate (OU=admin) |

## Usage

```bash
# Server
sudo worker-server --cert=certs/server.pem --key=certs/server-key.pem --ca=certs/ca.pem

# Client
worker-cli --cert=certs/client1.pem --key=certs/client1-key.pem --ca=certs/ca.pem start /bin/echo hello
```

## Adding New Clients

1. Create `clientN-csr.json` (copy from `client1-csr.json`, change CN)
2. Run: `cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client clientN-csr.json | cfssljson -bare clientN`

## Make Targets

```bash
make all          # Generate all certificates
make clean        # Remove generated certificates
make verify       # Show certificate details
make verify-chain # Verify certificate chain
```
