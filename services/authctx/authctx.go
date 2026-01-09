// Package authctx provides authentication context utilities for passing
// user information between HTTP handlers and gRPC services via metadata.
package authctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// Metadata keys for authentication context
const (
	// MetadataKeyUserID is the gRPC metadata key for the authenticated user ID
	MetadataKeyUserID = "x-user-id"

	// MetadataKeySwitchUser is the gRPC metadata key for switching to a different user (testing only)
	MetadataKeySwitchUser = "x-switch-user"
)

// UserIDFromContext extracts the authenticated user ID from the gRPC context metadata.
// Returns empty string if no user is authenticated.
func UserIDFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// Check for switch user first (only works if ENABLE_SWITCH_AUTH is set)
	if values := md.Get(MetadataKeySwitchUser); len(values) > 0 && values[0] != "" {
		return values[0]
	}

	// Get the actual user ID
	if values := md.Get(MetadataKeyUserID); len(values) > 0 {
		return values[0]
	}

	return ""
}

// UserIDToOutgoingContext adds the user ID to outgoing gRPC context metadata.
func UserIDToOutgoingContext(ctx context.Context, userID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, MetadataKeyUserID, userID)
}

// SwitchUserToOutgoingContext adds a switch-user header to outgoing gRPC context metadata.
// This is only effective when ENABLE_SWITCH_AUTH is set on the server.
func SwitchUserToOutgoingContext(ctx context.Context, switchToUserID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, MetadataKeySwitchUser, switchToUserID)
}

// IsAuthenticated returns true if there is an authenticated user in the context.
func IsAuthenticated(ctx context.Context) bool {
	return UserIDFromContext(ctx) != ""
}
