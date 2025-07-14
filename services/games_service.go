package services

import (
	"context"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// GamesServiceImpl implements the GamesService gRPC interface
type GamesServiceImpl struct {
	v1.UnimplementedGamesServiceServer
}

// NewGamesService creates a new GamesService implementation
func NewGamesService() *GamesServiceImpl {
	return &GamesServiceImpl{}
}

// ListGames returns all available games
func (s *GamesServiceImpl) ListGames(ctx context.Context, req *v1.ListGamesRequest) (resp *v1.ListGamesResponse, err error) {
	resp = &v1.ListGamesResponse{}
	return
}

// GetGame returns a specific game with metadata
func (s *GamesServiceImpl) GetGame(ctx context.Context, req *v1.GetGameRequest) (resp *v1.GetGameResponse, err error) {
	resp = &v1.GetGameResponse{}
	return
}
