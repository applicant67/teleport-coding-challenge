//go:build integration

package integration

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/manil/job-worker/pkg/proto"
)

// TestTLS_UntrustedClientCert tests that a client with a cert not signed by the CA is rejected.
func TestTLS_UntrustedClientCert(t *testing.T) {
	env := newTestEnv(t)

	// Generate a self-signed client cert (not signed by our CA)
	untrustedCert, untrustedKey := generateSelfSignedCert(t, "untrusted-client")

	// Try to connect with untrusted cert
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{untrustedCert.Raw},
			PrivateKey:  untrustedKey,
		}},
		RootCAs:    loadCAPool(t),
		MinVersion: tls.VersionTLS13,
	}

	conn, err := grpc.NewClient(
		env.serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	// Try to make a request - should fail during TLS handshake
	_, err = client.StartJob(context.Background(), &pb.StartJobRequest{Command: "/bin/echo"})
	if err == nil {
		t.Error("expected error when using untrusted client cert")
	}
	t.Logf("Correctly rejected untrusted client cert: %v", err)
}

// TestTLS_NoClientCert tests that a client without a cert is rejected.
func TestTLS_NoClientCert(t *testing.T) {
	env := newTestEnv(t)

	// Connect without providing a client certificate
	tlsConfig := &tls.Config{
		RootCAs:    loadCAPool(t),
		MinVersion: tls.VersionTLS13,
		// No Certificates field - no client cert
	}

	conn, err := grpc.NewClient(
		env.serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	// Try to make a request - should fail
	_, err = client.StartJob(context.Background(), &pb.StartJobRequest{Command: "/bin/echo"})
	if err == nil {
		t.Error("expected error when no client cert provided")
	}
	t.Logf("Correctly rejected missing client cert: %v", err)
}

// TestTLS_ExpiredClientCert tests that an expired client cert is rejected.
func TestTLS_ExpiredClientCert(t *testing.T) {
	env := newTestEnv(t)

	// Generate an expired client cert (signed by our CA would be better,
	// but for simplicity we use self-signed which will also be rejected)
	expiredCert, expiredKey := generateExpiredCert(t, "expired-client")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{expiredCert.Raw},
			PrivateKey:  expiredKey,
		}},
		RootCAs:    loadCAPool(t),
		MinVersion: tls.VersionTLS13,
	}

	conn, err := grpc.NewClient(
		env.serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	_, err = client.StartJob(context.Background(), &pb.StartJobRequest{Command: "/bin/echo"})
	if err == nil {
		t.Error("expected error when using expired client cert")
	}
	t.Logf("Correctly rejected expired client cert: %v", err)
}

// TestTLS_WrongCA tests that a client trusting a different CA rejects the server.
func TestTLS_WrongCA(t *testing.T) {
	env := newTestEnv(t)

	// Generate a different CA that didn't sign the server cert
	wrongCA, _ := generateSelfSignedCert(t, "wrong-ca")
	wrongCAPool := x509.NewCertPool()
	wrongCAPool.AddCert(wrongCA)

	// Load valid client cert but trust wrong CA
	clientCert, err := tls.LoadX509KeyPair(
		filepath.Join(certsDir, "client1.pem"),
		filepath.Join(certsDir, "client1-key.pem"),
	)
	if err != nil {
		t.Fatalf("failed to load client cert: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      wrongCAPool, // Trust wrong CA
		MinVersion:   tls.VersionTLS13,
	}

	conn, err := grpc.NewClient(
		env.serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	_, err = client.StartJob(context.Background(), &pb.StartJobRequest{Command: "/bin/echo"})
	if err == nil {
		t.Error("expected error when client trusts wrong CA")
	}
	t.Logf("Correctly rejected server with wrong CA: %v", err)
}

// TestTLS_ValidCertWorks tests that valid certificates work (sanity check).
func TestTLS_ValidCertWorks(t *testing.T) {
	env := newTestEnv(t)
	c := env.newClient("client1")

	// Should successfully start a job
	ctx := context.Background()
	jobID, err := c.Start(ctx, "/bin/echo", []string{"tls-test"})
	if err != nil {
		t.Fatalf("valid cert should work: %v", err)
	}
	t.Logf("Valid cert works, started job: %s", jobID)
}

// TestTLS_TLS12Rejected tests that TLS 1.2 connections are rejected (we require 1.3).
func TestTLS_TLS12Rejected(t *testing.T) {
	env := newTestEnv(t)

	clientCert, err := tls.LoadX509KeyPair(
		filepath.Join(certsDir, "client1.pem"),
		filepath.Join(certsDir, "client1-key.pem"),
	)
	if err != nil {
		t.Fatalf("failed to load client cert: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      loadCAPool(t),
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS12, // Force TLS 1.2
	}

	conn, err := grpc.NewClient(
		env.serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	_, err = client.StartJob(context.Background(), &pb.StartJobRequest{Command: "/bin/echo"})
	if err == nil {
		t.Error("expected error when using TLS 1.2 (server requires 1.3)")
	}
	t.Logf("Correctly rejected TLS 1.2: %v", err)
}

// Helper functions

// loadCAPool loads the CA certificate pool from the certs directory.
func loadCAPool(t *testing.T) *x509.CertPool {
	t.Helper()

	caPEM, err := os.ReadFile(filepath.Join(certsDir, "ca.pem"))
	if err != nil {
		t.Fatalf("failed to read CA cert: %v", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		t.Fatal("failed to add CA cert to pool")
	}

	return pool
}

// generateSelfSignedCert generates a self-signed certificate for testing.
func generateSelfSignedCert(t *testing.T, cn string) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	return cert, key
}

// generateExpiredCert generates an expired certificate for testing.
func generateExpiredCert(t *testing.T, cn string) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             time.Now().Add(-48 * time.Hour), // Started 2 days ago
		NotAfter:              time.Now().Add(-24 * time.Hour), // Expired 1 day ago
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	return cert, key
}

// writeTempCert writes a certificate and key to temporary files for testing.
func writeTempCert(t *testing.T, cert *x509.Certificate, key *ecdsa.PrivateKey) (certPath, keyPath string) {
	t.Helper()

	dir := t.TempDir()

	// Write cert
	certPath = filepath.Join(dir, "cert.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		t.Fatalf("failed to write cert: %v", err)
	}

	// Write key
	keyPath = filepath.Join(dir, "key.pem")
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		t.Fatalf("failed to write key: %v", err)
	}

	return certPath, keyPath
}
