//go:build !wasm
// +build !wasm

// Package authz provides authorization utilities for LilBattle services.
// It extracts user identity from gRPC context and validates access permissions.
package authz

import (
	"context"
	"fmt"
	"os"

	oagrpc "github.com/panyam/oneauth/grpc"
	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// isDevBypass reports whether the current process is running with
// DISABLE_API_AUTH=true. The gRPC interceptor already treats that env var
// as "auth optional at the transport"; extending the same bypass to
// handler-side RequireAuthenticated / CanSubmitMoves lets local-dev flows
// (make servelocal without a login, e2e replay harness) exercise the full
// move path without a real user. NEVER TRUE IN PRODUCTION — the env var
// is documented dev-only in web/server/dev_login.go's neighborhood.
func isDevBypass() bool {
	return os.Getenv("DISABLE_API_AUTH") == "true"
}

// Common authorization errors
var (
	ErrUnauthenticated = fmt.Errorf("authentication required")
	ErrForbidden       = fmt.Errorf("access denied")
	ErrNotOwner        = fmt.Errorf("you are not the owner of this resource")
	ErrNotPlayer       = fmt.Errorf("you are not a player in this game")
	ErrNotYourTurn     = fmt.Errorf("it is not your turn")
)

// GetUserIDFromContext extracts the authenticated user ID from gRPC context.
// Returns empty string if no user is authenticated.
func GetUserIDFromContext(ctx context.Context) string {
	return oagrpc.SubjectFromContext(ctx)
}

// RequireAuthenticated returns an error if no user is authenticated. In
// dev-bypass mode (DISABLE_API_AUTH=true) it returns "" with nil error so
// downstream authz calls can proceed on anonymous sessions — matches the
// interceptor's own "optional auth" behavior under the same env var.
func RequireAuthenticated(ctx context.Context) (string, error) {
	userID := GetUserIDFromContext(ctx)
	if userID == "" {
		if isDevBypass() {
			return "", nil
		}
		return "", ErrUnauthenticated
	}
	return userID, nil
}

// RequireOwnership checks if the authenticated user owns the resource.
// creatorID is the ID of the user who created/owns the resource.
// If creatorID is empty (orphaned resource), any authenticated user is allowed.
func RequireOwnership(ctx context.Context, creatorID string) error {
	_, err := RequireAuthenticated(ctx)
	if err != nil {
		return err
	}

	if creatorID == "" {
		// Orphaned resource — allow any authenticated user
		return nil
	}

	userID := GetUserIDFromContext(ctx)
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

// RequireGamePlayer checks if the authenticated user is a player in the game.
// Returns the player's ID (1-based) if they are a player.
func RequireGamePlayer(ctx context.Context, game *v1.Game) (int32, error) {
	userID, err := RequireAuthenticated(ctx)
	if err != nil {
		return 0, err
	}

	if game.Config == nil {
		return 0, ErrNotPlayer
	}

	for _, player := range game.Config.Players {
		if player.UserId == userID {
			return player.PlayerId, nil
		}
	}

	return 0, ErrNotPlayer
}

// RequireCurrentPlayer checks if it's the authenticated user's turn.
// Returns the player's ID if it's their turn.
func RequireCurrentPlayer(ctx context.Context, game *v1.Game, currentPlayer int32) (int32, error) {
	playerID, err := RequireGamePlayer(ctx, game)
	if err != nil {
		return 0, err
	}

	if playerID != currentPlayer {
		return playerID, ErrNotYourTurn
	}

	return playerID, nil
}

// CanSubmitMoves checks if user can submit moves to a game.
// User must be a player in the game AND it must be their turn.
//
// Dev-bypass: DISABLE_API_AUTH=true with no subject skips the check so the
// replay harness (and local-dev anonymous sessions) can drive both players
// from one connection. The bypass is transport-scoped by the interceptor
// upstream — flipping it in prod would take an explicit env-var flip on
// the deployed process, which is out-of-band from any request.
func CanSubmitMoves(ctx context.Context, game *v1.Game, currentPlayer int32) error {
	if isDevBypass() {
		return nil
	}
	_, err := RequireCurrentPlayer(ctx, game, currentPlayer)
	return err
}
