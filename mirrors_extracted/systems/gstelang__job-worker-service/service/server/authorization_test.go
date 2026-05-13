package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func TestIsAuthorized(t *testing.T) {
	tests := []struct {
		role       string
		methodName string
		expected   bool
	}{
		// Admin
		{"Admin", "StartJob", true},
		{"Admin", "QueryJob", true},
		{"Admin", "StopJob", true},
		{"Admin", "StreamLogs", true},
		// User
		{"User", "StartJob", false},
		{"User", "QueryJob", true},
		{"User", "StopJob", false},
		{"User", "StreamLogs", true},
		// Guest
		{"Guest", "QueryJob", false},
		{"Guest", "StreamLogs", false},
	}

	for _, test := range tests {
		t.Run(test.role+"_"+test.methodName, func(t *testing.T) {
			result := isAuthorized(test.role, test.methodName)
			if result != test.expected {
				t.Errorf("isAuthorized(%s, %s) = %v; expect authorized as %v", test.role, test.methodName, result, test.expected)
			}
		})
	}
}

func TestAuthorize(t *testing.T) {
	// Mock certificate with OID
	cert := &x509.Certificate{
		Extensions: []pkix.Extension{
			{
				Id:    []int{1, 3, 6, 1, 4, 1, 12345, 1, 1}, // OID for Admin
				Value: []byte("Admin"),
			},
		},
	}

	tlsInfo := credentials.TLSInfo{
		State: tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{cert},
		},
	}

	tests := []struct {
		methodName string
		expected   bool
	}{
		{"StartJob", true},
		{"StopJob", true},
		{"QueryJob", true},
		{"StreamLogs", true},
	}

	for _, test := range tests {
		t.Run(test.methodName, func(t *testing.T) {
			// Create a context with peer information
			p := &peer.Peer{
				AuthInfo: tlsInfo,
			}
			ctx := peer.NewContext(context.Background(), p)

			err := authorize(ctx, test.methodName)
			if (err == nil) != test.expected {
				t.Errorf("authorize(ctx, %s) => Got error as %v expect error to be %v", test.methodName, err, !test.expected)
			}
		})
	}
}

func TestAuthorize_NoPeer(t *testing.T) {
	ctx := context.Background() // No peer in context

	err := authorize(ctx, "StartJob")
	if err == nil {
		t.Error("should err if there's no peer")
	}
}

func TestAuthorize_UnknownRole(t *testing.T) {
	// Mock certificate with random OID
	cert := &x509.Certificate{
		Extensions: []pkix.Extension{
			{
				Id:    []int{1, 3, 6, 1, 4, 1, 12345, 1, 3},
				Value: []byte("Unknown"),
			},
		},
	}

	tlsInfo := credentials.TLSInfo{
		State: tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{cert},
		},
	}

	// Create a context with peer information
	p := &peer.Peer{
		AuthInfo: tlsInfo,
	}
	ctx := peer.NewContext(context.Background(), p)

	err := authorize(ctx, "StopJob")
	if err == nil {
		t.Error("StopJob should not be authorize for role unknown")
	}
}
