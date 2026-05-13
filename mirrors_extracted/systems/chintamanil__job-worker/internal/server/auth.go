package server

import (
	"context"
	"errors"
	"slices"

	"github.com/manil/job-worker/pkg/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ClientIdentity holds the authenticated client's identity extracted from mTLS certificate.
type ClientIdentity struct {
	UserID  string // From certificate CN (Common Name)
	IsAdmin bool   // True if certificate OU contains "admin"
}

// identityKeyType is a private type for context keys to avoid collisions.
// Using a pointer to an empty struct is more idiomatic than a string key.
type identityKeyType struct{}

// identityKey is the context key for storing client identity.
var identityKey = &identityKeyType{}

// ContextWithIdentity returns a new context with the client identity attached.
func ContextWithIdentity(ctx context.Context, identity *ClientIdentity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

// IdentityFromContext extracts the client identity from the context.
// Returns nil and false if no identity is present.
func IdentityFromContext(ctx context.Context) (*ClientIdentity, bool) {
	identity, ok := ctx.Value(identityKey).(*ClientIdentity)
	return identity, ok
}

// extractIdentityFromCert extracts client identity from the mTLS peer certificate.
// Returns an Unauthenticated error if no valid certificate is present.
func extractIdentityFromCert(ctx context.Context) (*ClientIdentity, error) {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil {
		return nil, status.Error(codes.Unauthenticated, "no peer in context")
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no TLS info in peer")
	}

	if len(tlsInfo.State.PeerCertificates) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no client certificate")
	}

	cert := tlsInfo.State.PeerCertificates[0]
	return &ClientIdentity{
		UserID:  cert.Subject.CommonName,
		IsAdmin: slices.Contains(cert.Subject.OrganizationalUnit, "admin"),
	}, nil
}

// mapError converts domain errors to gRPC status codes.
// Centralized here so handlers can return plain errors.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	// Check if already a gRPC status error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// Map domain errors to gRPC codes
	switch {
	case errors.Is(err, worker.ErrJobNotFound):
		return status.Error(codes.NotFound, "job not found")
	case errors.Is(err, worker.ErrJobAlreadyStarted):
		return status.Error(codes.FailedPrecondition, "job already started")
	case errors.Is(err, worker.ErrCgroupCreation):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, worker.ErrCgroupFDOpen):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, worker.ErrProcessStart):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, worker.ErrJobKill):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "deadline exceeded")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "canceled")
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// UnaryAuthInterceptor extracts client identity from mTLS certificate,
// attaches it to the request context, and maps errors to gRPC status codes.
func UnaryAuthInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	identity, err := extractIdentityFromCert(ctx)
	if err != nil {
		return nil, err
	}
	ctx = ContextWithIdentity(ctx, identity)

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// StreamAuthInterceptor extracts client identity from mTLS certificate,
// attaches it to the stream context, and maps errors to gRPC status codes.
func StreamAuthInterceptor(
	srv any,
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	identity, err := extractIdentityFromCert(ss.Context())
	if err != nil {
		return err
	}
	wrapped := &wrappedServerStream{
		ServerStream: ss,
		ctx:          ContextWithIdentity(ss.Context(), identity),
	}

	err = handler(srv, wrapped)
	if err != nil {
		return mapError(err)
	}
	return nil
}

// wrappedServerStream wraps grpc.ServerStream to override Context().
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
