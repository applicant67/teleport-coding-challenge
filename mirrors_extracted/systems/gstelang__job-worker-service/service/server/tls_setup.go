package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc/credentials"
)

func setupTLS() credentials.TransportCredentials {
	// Load server's certificate and private key
	certFile := "./certs/server/server.crt"
	certKey := "./certs/server/server.key"
	caFile := "./certs/ca/ca.crt"

	// Load CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load server's certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, certKey)
	if err != nil {
		panic(err)
	}

	// Set up TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS13, // Enforce TLS 1.3
		MaxVersion:   tls.VersionTLS13, // Enforce TLS 1.3
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return creds

}
