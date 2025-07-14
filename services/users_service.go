package services

import (
	"context"

	v1 "github.com/panyam/turnengine/games/weewar/gen/go/weewar/v1"
)

// UsersServiceImpl implements the UsersService gRPC interface
type UsersServiceImpl struct {
	v1.UnimplementedUsersServiceServer
}

// NewUsersService creates a new UsersService implementation
func NewUsersService() *UsersServiceImpl {
	return &UsersServiceImpl{}
}

// ListUsers returns all available users
func (s *UsersServiceImpl) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (resp *v1.ListUsersResponse, err error) {
	resp = &v1.ListUsersResponse{}
	return
}

// GetUser returns a specific user with metadata
func (s *UsersServiceImpl) GetUser(ctx context.Context, req *v1.GetUserRequest) (resp *v1.GetUserResponse, err error) {
	resp = &v1.GetUserResponse{}
	return
}
