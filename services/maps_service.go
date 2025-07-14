package services

import (
	"context"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// MapsServiceImpl implements the MapsService gRPC interface
type MapsServiceImpl struct {
	v1.UnimplementedMapsServiceServer
}

// NewMapsService creates a new MapsService implementation
func NewMapsService() *MapsServiceImpl {
	return &MapsServiceImpl{}
}

// ListMaps returns all available maps
func (s *MapsServiceImpl) ListMaps(ctx context.Context, req *v1.ListMapsRequest) (resp *v1.ListMapsResponse, err error) {
	resp = &v1.ListMapsResponse{}
	return
}

// GetMap returns a specific map with metadata
func (s *MapsServiceImpl) GetMap(ctx context.Context, req *v1.GetMapRequest) (resp *v1.GetMapResponse, err error) {
	resp = &v1.GetMapResponse{}
	return
}
