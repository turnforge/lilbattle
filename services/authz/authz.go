//go:build !wasm
// +build !wasm

// Package authz provides authorization utilities for WeeWar services.
// It extracts user identity from gRPC context and validates access permissions.
package authz

import (
	"context"
	"fmt"

	oagrpc "github.com/panyam/oneauth/grpc"
	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
)

// Common authorization errors
var (
	ErrUnauthenticated = fmt.Errorf("authentication required")
	ErrForbidden       = fmt.Errorf("access denied")
	ErrNotOwner        = fmt.Errorf("you are not the owner of this resource")
	ErrNotPlayer       = fmt.Errorf("you are not a player in this game")
)

// GetUserIDFromContext extracts the authenticated user ID from gRPC context.
// Returns empty string if no user is authenticated.
func GetUserIDFromContext(ctx context.Context) string {
	return oagrpc.UserIDFromContext(ctx)
}

// RequireAuthenticated returns an error if no user is authenticated.
func RequireAuthenticated(ctx context.Context) (string, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		return "", ErrUnauthenticated
	}
	return userID, nil
}

// RequireOwnership checks if the authenticated user owns the resource.
// creatorID is the ID of the user who created/owns the resource.
func RequireOwnership(ctx context.Context, creatorID string) error {
	userID, err := RequireAuthenticated(ctx)
	if err != nil {
		return err
	}

	if userID != creatorID {
		return ErrNotOwner
	}
	return nil
}

// CanModifyGame checks if user can modify game metadata (update/delete).
// Only the game creator can modify game metadata.
func CanModifyGame(ctx context.Context, game *v1.Game) error {
	return RequireOwnership(ctx, game.CreatorId)
}

// CanModifyWorld checks if user can modify world metadata (update/delete).
// Only the world creator can modify world metadata.
func CanModifyWorld(ctx context.Context, world *v1.World) error {
	return RequireOwnership(ctx, world.CreatorId)
}

// CanSubmitMoves checks if user can submit moves to a game.
// Currently checks if user is the game creator since GamePlayer doesn't
// have a user_id field yet. This works for:
// - Single player games (creator is the only player)
// - Hotseat games (all players share the same session)
//
// TODO: Add user_id to GamePlayer proto for proper multiplayer support
// where different users control different players.
func CanSubmitMoves(ctx context.Context, game *v1.Game) error {
	return RequireOwnership(ctx, game.CreatorId)
}
