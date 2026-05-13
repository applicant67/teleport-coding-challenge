#!/bin/bash -eu

pushd ssl

# generate CA certificates
openssl req -newkey rsa:2048 -nodes -x509 -days 365 -out ca.crt -keyout ca.key -subj "/CN=ca"

# generate server certificates
openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=server"
openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extfile <(echo subjectAltName = DNS:localhost)

# generate client certificates for user "alice"
openssl req -newkey rsa:2048 -nodes -keyout client-alice.key -out client-alice.csr -subj "/CN=alice"
openssl x509 -req -days 365 -sha256 -in client-alice.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client-alice.crt -extfile <(echo subjectAltName = DNS:localhost)

# generate client certificates for user "bob"
openssl req -newkey rsa:2048 -nodes -keyout client-bob.key -out client-bob.csr -subj "/CN=bob"
openssl x509 -req -days 365 -sha256 -in client-bob.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client-bob.crt -extfile <(echo subjectAltName = DNS:localhost)

# generate client certificates for user "unknown"
openssl req -newkey rsa:2048 -nodes -keyout client-unknown.key -out client-unknown.csr -subj "/CN=unknown"
openssl x509 -req -days 365 -sha256 -in client-unknown.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client-unknown.crt -extfile <(echo subjectAltName = DNS:localhost)

rm ca.srl server.csr client-alice.csr client-bob.csr client-unknown.csr || true

# openssl verify -verbose -CAfile ssl/ca.crt ssl/server.crt ssl/client-alice.crt ssl/client-invalid.crt
