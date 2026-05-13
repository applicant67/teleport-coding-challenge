package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc/credentials"
)

func setupTLS() credentials.TransportCredentials {

	// Load client certificates for admin
	// Hardcoding for admin
	// Ideally stored in a secret manager service on AWS/K8s serets or similar
	cert, err := tls.LoadX509KeyPair("./certs/client/admin/admin.crt", "./certs/client/admin/admin.key")
	if err != nil {
		log.Fatalf("failed to load client certificate: %v", err)
	}
	// Load CA certificate
	caCert, err := os.ReadFile("./certs/ca/ca.crt")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Set up TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})

	return creds
}
