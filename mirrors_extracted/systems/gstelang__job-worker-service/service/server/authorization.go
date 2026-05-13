package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// TODO: missed this. Move this logic to use UnaryClient/Streaming interceptors.

func isAuthorized(role, methodName string) bool {
	roleAccess := map[string][]string{
		"StartJob":   {"Admin"},
		"StopJob":    {"Admin"},
		"QueryJob":   {"Admin", "User"},
		"StreamLogs": {"Admin", "User"},
	}

	if allowedRoles, ok := roleAccess[methodName]; ok {
		for _, r := range allowedRoles {
			if role == r {
				return true
			}
		}
	}
	return false
}

func authorize(ctx context.Context, methodName string) error {
	// Hardcoded here but ideally I would store OIDs in AWS secrets manager if you're on AWS or similar.
	// These will be loaded dynamically in containerized env as ENV variables
	OIDRoleMapping := map[string]string{
		"1.3.6.1.4.1.12345.1.1": "Admin",
		"1.3.6.1.4.1.12345.1.2": "User",
	}

	peer, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("no peer found")
	}

	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return fmt.Errorf("unexpected peer transport credentials")
	}

	for _, cert := range tlsInfo.State.PeerCertificates {
		for _, ext := range cert.Extensions {
			oid := ext.Id.String()
			if role, exists := OIDRoleMapping[oid]; exists {
				if isAuthorized(role, methodName) {
					return nil
				}
				return fmt.Errorf("unauthorized: permission denied for role %s for action %s", role, methodName)
			}
		}
	}
	return fmt.Errorf("unknown role")
}
