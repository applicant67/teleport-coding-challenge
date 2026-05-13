#!/bin/bash

# Exit on error
set -e

# Determine the project root relative to this script (assuming scripts/generate_certs.sh)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CERTS_DIR="$PROJECT_ROOT/certs"

mkdir -p "$CERTS_DIR"
cd "$CERTS_DIR"

echo "Generating CA..."
# 1. Generate CA's Private Key and Self-Signed Certificate
# Added -sha256 for explicit signature algorithm strength
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -sha256 -days 365 -nodes -keyout ca.key -out ca.crt -subj "/C=US/ST=State/L=City/O=Teleport/OU=Eng/CN=TeleportCA"

echo "Generating Server Certificate..."
# 2. Generate Server Private Key and CSR
openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes -keyout server.key -out server.csr -subj "/C=US/ST=State/L=City/O=Teleport/OU=Eng/CN=localhost"

# 3. Sign Server Certificate with CA (adding SAN for localhost)
# Go requires SANs (Subject Alternative Names) for IP/DNS verification.
echo "subjectAltName=DNS:localhost,IP:127.0.0.1" > server-ext.cnf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extfile server-ext.cnf

echo "Generating Client 'Alice' Certificate..."
# 4. Generate Client 'Alice' Private Key and CSR
openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes -keyout alice.key -out alice.csr -subj "/C=US/ST=State/L=City/O=Teleport/OU=Eng/CN=alice"

# 5. Sign Client 'Alice' Certificate with CA
openssl x509 -req -in alice.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out alice.crt -days 365 -sha256

echo "Generating Client 'Bob' Certificate..."
# 6. Generate Client 'Bob' Private Key and CSR
openssl req -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes -keyout bob.key -out bob.csr -subj "/C=US/ST=State/L=City/O=Teleport/OU=Eng/CN=bob"

# 7. Sign Client 'Bob' Certificate with CA
openssl x509 -req -in bob.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out bob.crt -days 365 -sha256

# Cleanup CSRs and config
rm *.csr server-ext.cnf

echo "Certificates generated in $CERTS_DIR:"
ls -l