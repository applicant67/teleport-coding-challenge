#!/bin/bash
# Generate TLS certificates for job-manager service.
#
# Usage: ./gen-certs.sh [certs_dir] [dns_name] [ip_addr]
#
# Examples:
#   ./gen-certs.sh                                    # defaults: certs/, localhost, 127.0.0.1
#   ./gen-certs.sh certs myserver.local               # custom DNS
#   ./gen-certs.sh certs myserver.local 192.168.1.10  # custom DNS + IP
#
# Creates:
#   certs/ca.crt, certs/ca.key         - Certificate Authority
#   certs/server.crt, certs/server.key - Server certificate
#   certs/admin.crt, certs/admin.key   - Admin client certificate
#   certs/read.crt, certs/read.key     - Read-only client certificate

set -e

CERTS_DIR="${1:-certs}"
DNS_NAME="${2:-localhost}"
IP_ADDR="${3:-127.0.0.1}"
DAYS=365

mkdir -p "$CERTS_DIR"
cd "$CERTS_DIR"

echo "==> Generating CA..."
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days $DAYS -key ca.key -out ca.crt \
    -subj "/CN=Job Manager CA/O=JobManager"

echo "==> Generating server certificate (DNS:${DNS_NAME}, IP:${IP_ADDR})..."
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr \
    -subj "/CN=${DNS_NAME}/O=JobManager"

cat > server.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage=digitalSignature,keyEncipherment
extendedKeyUsage=serverAuth
subjectAltName=DNS:${DNS_NAME},IP:${IP_ADDR}
EOF

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days $DAYS -extfile server.ext
rm server.csr server.ext

echo "==> Generating admin client certificate..."
openssl genrsa -out admin.key 4096
openssl req -new -key admin.key -out admin.csr \
    -subj "/CN=admin/O=JobManager/OU=Admin"

cat > client.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage=digitalSignature
extendedKeyUsage=clientAuth
EOF

openssl x509 -req -in admin.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out admin.crt -days $DAYS -extfile client.ext
rm admin.csr

echo "==> Generating read-only client certificate..."
openssl genrsa -out read.key 4096
openssl req -new -key read.key -out read.csr \
    -subj "/CN=reader/O=JobManager/OU=Read"

openssl x509 -req -in read.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out read.crt -days $DAYS -extfile client.ext
rm read.csr client.ext

rm -f ca.srl

echo "==> Certificates generated in $CERTS_DIR/"
ls -la