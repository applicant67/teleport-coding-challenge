#!/bin/bash
# file: generate-certs.sh

# Create necessary directories
mkdir -p ./certs/ca ./certs/server ./certs/client/user ./certs/client/admin ./certs/cnf

# Generate CA key and cert
echo "Generate the Certificate Authority (CA) Key and Certificate..."
openssl genrsa -out ./certs/ca/ca.key 4096
openssl req -x509 -new -nodes -sha256 -days 365 -config ./certs/cnf/ca.cnf -key ./certs/ca/ca.key -out ./certs/ca/ca.crt

# Generate server key and cert
echo "Generate the Server Key and CSR..."
openssl genrsa -out ./certs/server/server.key 4096
openssl req -new -config ./certs/cnf/server.cnf -key ./certs/server/server.key -out ./certs/server/server.csr
echo "Sign the Server CSR with the CA..."
openssl x509 -req -in ./certs/server/server.csr -CA ./certs/ca/ca.crt -CAkey ./certs/ca/ca.key -CAcreateserial -out ./certs/server/server.crt -days 365 -sha256 -extfile ./certs/cnf/server.cnf -extensions v3_req

# Client (User role)
echo "Client user role cert generation..."
echo "# 1. First, generate the key..."
openssl genrsa -out ./certs/client/user/user.key 4096
echo "# 2. Then, generate the CSR..."
openssl req -new -config ./certs/cnf/user.cnf -key ./certs/client/user/user.key -out ./certs/client/user/user.csr 
echo "# 3. Sign the CSR with the CA to generate the certificate..."
openssl x509 -req -in ./certs/client/user/user.csr -CA ./certs/ca/ca.crt -CAkey ./certs/ca/ca.key -CAcreateserial -out ./certs/client/user/user.crt -days 45 -sha256 -extfile ./certs/cnf/user.cnf -extensions v3_req

# Client (Admin role)
echo "Client admin role cert generation..."
echo "# 1. First, generate the key..."
openssl genrsa -out ./certs/client/admin/admin.key 4096
echo "# 2. Then, generate the CSR..."
openssl req -new -config ./certs/cnf/admin.cnf -key ./certs/client/admin/admin.key -out ./certs/client/admin/admin.csr 
echo "# 3. Sign the CSR with the CA to generate the certificate..."
openssl x509 -req -in ./certs/client/admin/admin.csr -CA ./certs/ca/ca.crt -CAkey ./certs/ca/ca.key -CAcreateserial -out ./certs/client/admin/admin.crt -days 45 -sha256 -extfile ./certs/cnf/admin.cnf -extensions v3_req

# View certificate details to verify correctness
# openssl x509 -in ./certs/ca/ca.crt -text -noout
# openssl x509 -in ./certs/client/admin/admin.crt -text -noout
# openssl x509 -in ./certs/client/user/user.crt -text -noout
# openssl x509 -in ./certs/server/server.crt -text -noout
