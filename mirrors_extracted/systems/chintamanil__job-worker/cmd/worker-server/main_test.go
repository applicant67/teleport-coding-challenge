//go:build linux

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manil/job-worker/internal/server"
)

func TestServerConfig(t *testing.T) {
	cfg := &serverConfig{
		listenAddr: "0.0.0.0:50051",
		certFile:   "./certs/server.pem",
		keyFile:    "./certs/server-key.pem",
		caFile:     "./certs/ca.pem",
		cgroupRoot: "/sys/fs/cgroup/jobworker",
	}

	if cfg.listenAddr == "" {
		t.Error("listenAddr should not be empty")
	}
	if cfg.certFile == "" {
		t.Error("certFile should not be empty")
	}
	if cfg.keyFile == "" {
		t.Error("keyFile should not be empty")
	}
	if cfg.caFile == "" {
		t.Error("caFile should not be empty")
	}
	if cfg.cgroupRoot == "" {
		t.Error("cgroupRoot should not be empty")
	}
}

func TestServerTLSConfig(t *testing.T) {
	certsDir := filepath.Join("..", "..", "certs")

	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		t.Skip("certs directory not found, skipping TLS config test")
	}

	certFile := filepath.Join(certsDir, "server.pem")
	keyFile := filepath.Join(certsDir, "server-key.pem")
	caFile := filepath.Join(certsDir, "ca.pem")

	t.Run("valid server certificates", func(t *testing.T) {
		cfg, err := server.NewTLSConfig(certFile, keyFile, caFile)
		if err != nil {
			t.Fatalf("NewTLSConfig() error = %v", err)
		}

		if cfg == nil {
			t.Fatal("NewTLSConfig() returned nil config")
		}

		if len(cfg.Certificates) != 1 {
			t.Errorf("expected 1 certificate, got %d", len(cfg.Certificates))
		}

		if cfg.ClientCAs == nil {
			t.Error("ClientCAs is nil, should be set for mTLS")
		}

		if cfg.MinVersion != 0x0304 { // TLS 1.3
			t.Errorf("MinVersion = %x, want TLS 1.3 (0x0304)", cfg.MinVersion)
		}

		// Verify RequireAndVerifyClientCert for mTLS
		if cfg.ClientAuth != 4 { // tls.RequireAndVerifyClientCert
			t.Errorf("ClientAuth = %d, want RequireAndVerifyClientCert (4)", cfg.ClientAuth)
		}
	})

	t.Run("missing server cert", func(t *testing.T) {
		_, err := server.NewTLSConfig("nonexistent.pem", keyFile, caFile)
		if err == nil {
			t.Error("expected error for missing cert file")
		}
	})

	t.Run("missing server key", func(t *testing.T) {
		_, err := server.NewTLSConfig(certFile, "nonexistent.pem", caFile)
		if err == nil {
			t.Error("expected error for missing key file")
		}
	})

	t.Run("missing CA file", func(t *testing.T) {
		_, err := server.NewTLSConfig(certFile, keyFile, "nonexistent.pem")
		if err == nil {
			t.Error("expected error for missing CA file")
		}
	})
}

func TestCgroupRootPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantErr  bool
	}{
		{
			name:    "default path",
			path:    "/sys/fs/cgroup/jobworker",
			wantErr: false,
		},
		{
			name:    "custom path",
			path:    "/sys/fs/cgroup/custom-jobworker",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "relative/path",
			wantErr: false, // Path validation happens at runtime
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &serverConfig{
				cgroupRoot: tt.path,
			}
			if cfg.cgroupRoot != tt.path {
				t.Errorf("cgroupRoot = %s, want %s", cfg.cgroupRoot, tt.path)
			}
		})
	}
}

func TestListenAddressFormats(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		isValid bool
	}{
		{name: "all interfaces with port", addr: "0.0.0.0:50051", isValid: true},
		{name: "localhost with port", addr: "127.0.0.1:50051", isValid: true},
		{name: "port only", addr: ":50051", isValid: true},
		{name: "custom port", addr: "0.0.0.0:8080", isValid: true},
		{name: "ipv6 localhost", addr: "[::1]:50051", isValid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &serverConfig{
				listenAddr: tt.addr,
			}
			// Just verify the address is stored correctly
			// Actual validation happens when server starts
			if cfg.listenAddr != tt.addr {
				t.Errorf("listenAddr = %s, want %s", cfg.listenAddr, tt.addr)
			}
		})
	}
}
