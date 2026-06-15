//go:build !wasm
// +build !wasm

package server

import (
	"context"
	"testing"

	"github.com/panyam/oneauth/core"
	oagrpc "github.com/panyam/oneauth/grpc"
)

func TestInjectAuthMetadata_PropagatesUserID(t *testing.T) {
	// Simulate what the HTTP auth middleware does: set user ID in context
	ctx := core.SetSubjectInContext(context.Background(), "user-123")

	// Run the Connect adapter's metadata injection
	ctx = injectAuthMetadata(ctx)

	// Verify the service can read the user ID via gRPC incoming metadata
	got := oagrpc.SubjectFromContext(ctx)
	if got != "user-123" {
		t.Errorf("UserIDFromContext = %q, want %q", got, "user-123")
	}
}

func TestInjectAuthMetadata_EmptyUserID(t *testing.T) {
	// No user ID set in context
	ctx := context.Background()

	ctx = injectAuthMetadata(ctx)

	got := oagrpc.SubjectFromContext(ctx)
	if got != "" {
		t.Errorf("UserIDFromContext = %q, want empty string", got)
	}
}
