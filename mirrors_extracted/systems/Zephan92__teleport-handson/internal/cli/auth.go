// Package cli implements the command-line interface for the Job Worker Service.
package cli

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	pb "github.com/teleport-handson/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// withClient handles the gRPC connection lifecycle for a command.
// It establishes a connection using the configuration, executes the provided function,
// and ensures the connection is closed.
func (c *Config) withClient(fn func(pb.JobWorkerClient) error) error {
	client, conn, err := c.Connect()
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Close()
	}
	return fn(client)
}

// dial establishes a gRPC connection using mTLS.
// It loads the TLS configuration and dials the server.
func (c *Config) dial() (pb.JobWorkerClient, *grpc.ClientConn, error) {
	tlsConfig, err := c.loadTLSConfig()
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.NewClient(c.ServerAddr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}

	return pb.NewJobWorkerClient(conn), conn, nil
}

// loadTLSConfig loads the TLS configuration for mTLS.
// It reads the CA certificate, client certificate, and client key from the paths specified in Config.
func (c *Config) loadTLSConfig() (*tls.Config, error) {
	pemServerCA, err := os.ReadFile(c.CAPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA: %w", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to append CA")
	}

	clientCert, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert/key: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}, nil
}
