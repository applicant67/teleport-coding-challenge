// Package auth provides mTLS authentication and role-based authorization.
package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

var (
	ErrNoPeerInfo       = errors.New("no peer info in context")
	ErrNoTLSInfo        = errors.New("no TLS info in peer")
	ErrNoCertificate    = errors.New("no client certificate")
	ErrInvalidRole      = errors.New("invalid role in certificate")
	ErrPermissionDenied = errors.New("permission denied")
)

// Role represents a client's authorization level.
type Role string

const (
	RoleAdmin Role = "admin" // Full access
	RoleRead  Role = "read"  // Read-only access
)

func (r Role) String() string {
	return string(r)
}

// ParseRole parses a role string from certificate OU.
func ParseRole(s string) (Role, error) {
	switch s {
	case "Admin", "admin":
		return RoleAdmin, nil
	case "Read", "read":
		return RoleRead, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidRole, s)
	}
}

// Permission represents an operation that requires authorization.
type Permission string

const (
	PermStart        Permission = "start"
	PermStop         Permission = "stop"
	PermStatus       Permission = "status"
	PermStreamOutput Permission = "stream_output"
)

// rolePermissions defines what each role can do.
var rolePermissions = map[Role]map[Permission]bool{
	RoleAdmin: {
		PermStart:        true,
		PermStop:         true,
		PermStatus:       true,
		PermStreamOutput: true,
	},
	RoleRead: {
		PermStart:        false,
		PermStop:         false,
		PermStatus:       true,
		PermStreamOutput: true,
	},
}

// HasPermission checks if a role can perform an operation.
func (r Role) HasPermission(p Permission) bool {
	perms, ok := rolePermissions[r]
	if !ok {
		return false
	}
	return perms[p]
}

// Identity represents an authenticated client.
type Identity struct {
	Role        Role
	CommonName  string
	Certificate *x509.Certificate
}

// ExtractIdentity extracts client identity from gRPC context.
// The identity comes from the mTLS client certificate.
func ExtractIdentity(ctx context.Context) (*Identity, error) {
	// Get peer info from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrNoPeerInfo
	}

	// Get TLS info
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, ErrNoTLSInfo
	}

	// Get verified certificate
	if len(tlsInfo.State.VerifiedChains) == 0 ||
		len(tlsInfo.State.VerifiedChains[0]) == 0 {
		return nil, ErrNoCertificate
	}

	cert := tlsInfo.State.VerifiedChains[0][0]

	// Extract role from OrganizationalUnit
	if len(cert.Subject.OrganizationalUnit) != 1 {
		return nil, fmt.Errorf("%w: need exactly 1 OU, got %d",
			ErrInvalidRole, len(cert.Subject.OrganizationalUnit))
	}

	role, err := ParseRole(cert.Subject.OrganizationalUnit[0])
	if err != nil {
		return nil, err
	}

	return &Identity{
		Role:        role,
		CommonName:  cert.Subject.CommonName,
		Certificate: cert,
	}, nil
}

// Authorize checks if an identity can perform an operation.
func Authorize(identity *Identity, perm Permission) error {
	if !identity.Role.HasPermission(perm) {
		return fmt.Errorf("%w: role %q cannot %q",
			ErrPermissionDenied, identity.Role, perm)
	}
	return nil
}

// TLSConfig holds certificate paths.
type TLSConfig struct {
	CertFile string // Server/client certificate
	KeyFile  string // Private key
	CAFile   string // CA certificate
}

// LoadServerTLSConfig creates TLS config for the server.
// Enforces mTLS with client certificate verification.
func LoadServerTLSConfig(cfg TLSConfig) (*tls.Config, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load server cert: %w", err)
	}

	// Load CA for client verification
	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// LoadClientTLSConfig creates TLS config for clients.
func LoadClientTLSConfig(cfg TLSConfig) (*tls.Config, error) {
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}

	// Load CA for server verification
	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to parse CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
