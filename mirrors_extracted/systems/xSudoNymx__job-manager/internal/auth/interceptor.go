package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// contextKey is used for storing identity in context.
type contextKey string

const identityKey contextKey = "identity"

// methodPermissions maps gRPC methods to required permissions.
var methodPermissions = map[string]Permission{
	"/job.v1.JobService/Start":        PermStart,
	"/job.v1.JobService/Stop":         PermStop,
	"/job.v1.JobService/Status":       PermStatus,
	"/job.v1.JobService/StreamOutput": PermStreamOutput,
}

// UnaryInterceptor creates an interceptor for unary RPCs.
func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract identity from TLS certificate
		identity, err := ExtractIdentity(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated,
				"authentication failed")
		}

		// Check authorization
		perm, ok := methodPermissions[info.FullMethod]
		if !ok {
			return nil, status.Errorf(codes.Unimplemented,
				"unknown method")
		}

		if err := Authorize(identity, perm); err != nil {
			return nil, status.Errorf(codes.PermissionDenied,
				"authorization failed")
		}

		// Store identity in context
		ctx = context.WithValue(ctx, identityKey, identity)

		return handler(ctx, req)
	}
}

// StreamInterceptor creates an interceptor for streaming RPCs.
func StreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		// Extract identity
		identity, err := ExtractIdentity(ctx)
		if err != nil {
			return status.Errorf(codes.Unauthenticated,
				"authentication failed")
		}

		// Check authorization
		perm, ok := methodPermissions[info.FullMethod]
		if !ok {
			return status.Errorf(codes.Unimplemented,
				"unknown method")
		}

		if err := Authorize(identity, perm); err != nil {
			return status.Errorf(codes.PermissionDenied,
				"authorization failed")
		}

		// Wrap stream with identity in context
		wrapped := &wrappedStream{
			ServerStream: ss,
			ctx:          context.WithValue(ctx, identityKey, identity),
		}

		return handler(srv, wrapped)
	}
}

// wrappedStream wraps ServerStream to inject custom context.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// IdentityFromContext retrieves identity from context.
func IdentityFromContext(ctx context.Context) *Identity {
	id, _ := ctx.Value(identityKey).(*Identity)
	return id
}
