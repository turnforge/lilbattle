//go:build !wasm
// +build !wasm

package tests

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// TestUserID is the default user ID used for testing
const TestUserID = "test-user-1"

// ContextWithUserID creates a context with the user ID set in gRPC metadata.
// This simulates what the auth interceptor does in production.
// Uses "x-user-id" which is oagrpc.DefaultMetadataKeyUserID from oneauth.
func ContextWithUserID(userID string) context.Context {
	md := metadata.Pairs("x-user-id", userID)
	return metadata.NewIncomingContext(context.Background(), md)
}

// AuthenticatedContext returns a context authenticated with the default test user.
func AuthenticatedContext() context.Context {
	return ContextWithUserID(TestUserID)
}
