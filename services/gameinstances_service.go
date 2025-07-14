package services

import (
	"context"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// GameInstancesServiceImpl implements the GameInstancesService gRPC interface
type GameInstancesServiceImpl struct {
	v1.UnimplementedGameInstancesServiceServer
}

// NewGameInstancesService creates a new GameInstancesService implementation
func NewGameInstancesService() *GameInstancesServiceImpl {
	return &GameInstancesServiceImpl{}
}

// ListGameInstances returns all available gameinstances
func (s *GameInstancesServiceImpl) ListGameInstances(ctx context.Context, req *v1.ListGameInstancesRequest) (resp *v1.ListGameInstancesResponse, err error) {
	resp = &v1.ListGameInstancesResponse{}
	return
}

// GetGameInstance returns a specific gameinstance with metadata
func (s *GameInstancesServiceImpl) GetGameInstance(ctx context.Context, req *v1.GetGameInstanceRequest) (resp *v1.GetGameInstanceResponse, err error) {
	resp = &v1.GetGameInstanceResponse{}
	return
}
