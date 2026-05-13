package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadTLSConfig(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath := generateTestCert(t, tmpDir)

	cfg := &Config{
		CAPath:   certPath, // Use the same cert as CA for self-signed test
		CertPath: certPath,
		KeyPath:  keyPath,
	}

	tlsConfig, err := cfg.loadTLSConfig()
	if err != nil {
		t.Fatalf("loadTLSConfig failed: %v", err)
	}

	if len(tlsConfig.Certificates) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(tlsConfig.Certificates))
	}
	if tlsConfig.RootCAs == nil {
		t.Error("Expected RootCAs to be initialized")
	}
}

func TestLoadTLSConfig_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	certPath, keyPath := generateTestCert(t, tmpDir)

	tests := []struct {
		name    string
		cfg     *Config
		wantErr string
	}{
		{
			name: "MissingCA",
			cfg: &Config{
				CAPath:   filepath.Join(tmpDir, "missing-ca.crt"),
				CertPath: certPath,
				KeyPath:  keyPath,
			},
			wantErr: "failed to load CA",
		},
		{
			name: "InvalidCA",
			cfg: func() *Config {
				badCA := filepath.Join(tmpDir, "bad-ca.crt")
				_ = os.WriteFile(badCA, []byte("not a certificate"), 0644)
				return &Config{
					CAPath:   badCA,
					CertPath: certPath,
					KeyPath:  keyPath,
				}
			}(),
			wantErr: "failed to append CA",
		},
		{
			name: "MissingCert",
			cfg: &Config{
				CAPath:   certPath,
				CertPath: filepath.Join(tmpDir, "missing.crt"),
				KeyPath:  keyPath,
			},
			wantErr: "failed to load client cert/key",
		},
		{
			name: "InvalidCertContent",
			cfg: func() *Config {
				badCert := filepath.Join(tmpDir, "bad.crt")
				_ = os.WriteFile(badCert, []byte("bad cert"), 0644)
				return &Config{
					CAPath:   certPath,
					CertPath: badCert,
					KeyPath:  keyPath,
				}
			}(),
			wantErr: "failed to load client cert/key",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.cfg.loadTLSConfig()
			if err == nil {
				t.Error("loadTLSConfig() expected error, got nil")
			} else if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("loadTLSConfig() error = %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}

func generateTestCert(t *testing.T, dir string) (string, string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Test"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}

	certPath := filepath.Join(dir, "cert.pem")
	certOut, err := os.Create(certPath)
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyPath := filepath.Join(dir, "key.pem")
	keyOut, err := os.Create(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	keyOut.Close()

	return certPath, keyPath
}
