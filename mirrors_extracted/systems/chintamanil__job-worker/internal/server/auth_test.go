//go:build linux

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/manil/job-worker/pkg/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func TestIdentityFromContext_NotSet(t *testing.T) {
	ctx := context.Background()
	identity, ok := IdentityFromContext(ctx)
	if ok {
		t.Error("expected ok=false when identity not set")
	}
	if identity != nil {
		t.Error("expected nil identity when not set")
	}
}

func TestIdentityFromContext_Set(t *testing.T) {
	ctx := context.Background()
	expected := &ClientIdentity{UserID: "client1", IsAdmin: false}
	ctx = ContextWithIdentity(ctx, expected)

	identity, ok := IdentityFromContext(ctx)
	if !ok {
		t.Error("expected ok=true when identity is set")
	}
	if identity == nil {
		t.Fatal("expected non-nil identity")
	}
	if identity.UserID != "client1" {
		t.Errorf("expected UserID=client1, got %s", identity.UserID)
	}
	if identity.IsAdmin {
		t.Error("expected IsAdmin=false")
	}
}

func TestIdentityFromContext_Admin(t *testing.T) {
	ctx := context.Background()
	expected := &ClientIdentity{UserID: "admin", IsAdmin: true}
	ctx = ContextWithIdentity(ctx, expected)

	identity, ok := IdentityFromContext(ctx)
	if !ok {
		t.Error("expected ok=true")
	}
	if !identity.IsAdmin {
		t.Error("expected IsAdmin=true for admin user")
	}
}

// createTestCert creates a test certificate for testing.
func createTestCert(t *testing.T, cn string, ou []string) *x509.Certificate {
	t.Helper()
	return &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         cn,
			OrganizationalUnit: ou,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),
	}
}

// createPeerContext creates a context with TLS peer info for testing.
func createPeerContext(t *testing.T, cert *x509.Certificate) context.Context {
	t.Helper()
	tlsInfo := credentials.TLSInfo{
		State: tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{cert},
		},
	}
	p := &peer.Peer{
		Addr:     &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
		AuthInfo: tlsInfo,
	}
	return peer.NewContext(context.Background(), p)
}

func TestExtractIdentityFromCert_RegularUser(t *testing.T) {
	cert := createTestCert(t, "client1", []string{"user"})
	ctx := createPeerContext(t, cert)

	identity, err := extractIdentityFromCert(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.UserID != "client1" {
		t.Errorf("expected UserID=client1, got %s", identity.UserID)
	}
	if identity.IsAdmin {
		t.Error("expected IsAdmin=false for regular user")
	}
}

func TestExtractIdentityFromCert_AdminUser(t *testing.T) {
	cert := createTestCert(t, "admin", []string{"admin"})
	ctx := createPeerContext(t, cert)

	identity, err := extractIdentityFromCert(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.UserID != "admin" {
		t.Errorf("expected UserID=admin, got %s", identity.UserID)
	}
	if !identity.IsAdmin {
		t.Error("expected IsAdmin=true for admin user")
	}
}

func TestExtractIdentityFromCert_NoPeer(t *testing.T) {
	ctx := context.Background()
	_, err := extractIdentityFromCert(ctx)
	if err == nil {
		t.Error("expected error when no peer in context")
	}
}

func TestExtractIdentityFromCert_NoCertificates(t *testing.T) {
	tlsInfo := credentials.TLSInfo{
		State: tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{},
		},
	}
	p := &peer.Peer{
		Addr:     &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
		AuthInfo: tlsInfo,
	}
	ctx := peer.NewContext(context.Background(), p)

	_, err := extractIdentityFromCert(ctx)
	if err == nil {
		t.Error("expected error when no certificates")
	}
}

func TestUnaryAuthInterceptor_Success(t *testing.T) {
	cert := createTestCert(t, "client1", []string{"user"})
	ctx := createPeerContext(t, cert)

	var capturedCtx context.Context
	handler := func(ctx context.Context, req any) (any, error) {
		capturedCtx = ctx
		return "response", nil
	}

	resp, err := UnaryAuthInterceptor(ctx, "request", &grpc.UnaryServerInfo{}, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "response" {
		t.Errorf("expected response, got %v", resp)
	}

	// Verify identity was attached to context
	identity, ok := IdentityFromContext(capturedCtx)
	if !ok {
		t.Fatal("expected identity in context")
	}
	if identity.UserID != "client1" {
		t.Errorf("expected UserID=client1, got %s", identity.UserID)
	}
}

func TestUnaryAuthInterceptor_NoPeer(t *testing.T) {
	ctx := context.Background()
	handler := func(ctx context.Context, req any) (any, error) {
		t.Error("handler should not be called")
		return nil, nil
	}

	_, err := UnaryAuthInterceptor(ctx, "request", &grpc.UnaryServerInfo{}, handler)
	if err == nil {
		t.Error("expected error when no peer")
	}
}

func TestUnaryAuthInterceptor_HandlerError(t *testing.T) {
	cert := createTestCert(t, "client1", []string{"user"})
	ctx := createPeerContext(t, cert)

	handlerErr := errors.New("handler error")
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, handlerErr
	}

	_, err := UnaryAuthInterceptor(ctx, "request", &grpc.UnaryServerInfo{}, handler)
	if err == nil {
		t.Error("expected error when handler fails")
	}
}

// mockServerStream implements grpc.ServerStream for testing.
type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func TestStreamAuthInterceptor_Success(t *testing.T) {
	cert := createTestCert(t, "client1", []string{"user"})
	ctx := createPeerContext(t, cert)
	stream := &mockServerStream{ctx: ctx}

	var capturedIdentity *ClientIdentity
	handler := func(srv any, ss grpc.ServerStream) error {
		var ok bool
		capturedIdentity, ok = IdentityFromContext(ss.Context())
		if !ok {
			return errors.New("no identity")
		}
		return nil
	}

	err := StreamAuthInterceptor(nil, stream, &grpc.StreamServerInfo{}, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedIdentity == nil {
		t.Fatal("expected identity to be set")
	}
	if capturedIdentity.UserID != "client1" {
		t.Errorf("expected UserID=client1, got %s", capturedIdentity.UserID)
	}
}

func TestStreamAuthInterceptor_NoPeer(t *testing.T) {
	ctx := context.Background()
	stream := &mockServerStream{ctx: ctx}

	handler := func(srv any, ss grpc.ServerStream) error {
		t.Error("handler should not be called")
		return nil
	}

	err := StreamAuthInterceptor(nil, stream, &grpc.StreamServerInfo{}, handler)
	if err == nil {
		t.Error("expected error when no peer")
	}
}

func TestStreamAuthInterceptor_HandlerError(t *testing.T) {
	cert := createTestCert(t, "client1", []string{"user"})
	ctx := createPeerContext(t, cert)
	stream := &mockServerStream{ctx: ctx}

	handlerErr := errors.New("handler error")
	handler := func(srv any, ss grpc.ServerStream) error {
		return handlerErr
	}

	err := StreamAuthInterceptor(nil, stream, &grpc.StreamServerInfo{}, handler)
	if err == nil {
		t.Error("expected error when handler fails")
	}
}

func TestExtractIdentityFromCert_NoTLSInfo(t *testing.T) {
	// Create peer with non-TLS auth info
	p := &peer.Peer{
		Addr:     &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345},
		AuthInfo: nil, // No auth info
	}
	ctx := peer.NewContext(context.Background(), p)

	_, err := extractIdentityFromCert(ctx)
	if err == nil {
		t.Error("expected error when no TLS info")
	}
}

func TestMapError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode codes.Code
	}{
		{
			name:         "nil error",
			err:          nil,
			expectedCode: codes.OK,
		},
		{
			name:         "job not found",
			err:          worker.ErrJobNotFound,
			expectedCode: codes.NotFound,
		},
		{
			name:         "job already started",
			err:          worker.ErrJobAlreadyStarted,
			expectedCode: codes.FailedPrecondition,
		},
		{
			name:         "cgroup creation error",
			err:          fmt.Errorf("test: %w", worker.ErrCgroupCreation),
			expectedCode: codes.Internal,
		},
		{
			name:         "cgroup fd open error",
			err:          fmt.Errorf("test: %w", worker.ErrCgroupFDOpen),
			expectedCode: codes.Internal,
		},
		{
			name:         "process start error",
			err:          fmt.Errorf("test: %w", worker.ErrProcessStart),
			expectedCode: codes.Internal,
		},
		{
			name:         "job kill error",
			err:          fmt.Errorf("test: %w", worker.ErrJobKill),
			expectedCode: codes.Internal,
		},
		{
			name:         "deadline exceeded",
			err:          context.DeadlineExceeded,
			expectedCode: codes.DeadlineExceeded,
		},
		{
			name:         "context canceled",
			err:          context.Canceled,
			expectedCode: codes.Canceled,
		},
		{
			name:         "unknown error",
			err:          errors.New("unknown error"),
			expectedCode: codes.Internal,
		},
		{
			name:         "already grpc status error",
			err:          status.Error(codes.PermissionDenied, "denied"),
			expectedCode: codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapError(tt.err)
			if tt.err == nil {
				if result != nil {
					t.Errorf("expected nil for nil input, got %v", result)
				}
				return
			}
			gotCode := status.Code(result)
			if gotCode != tt.expectedCode {
				t.Errorf("expected code %v, got %v", tt.expectedCode, gotCode)
			}
		})
	}
}
