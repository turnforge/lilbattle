package services

import (
	"context"
	"sync"
	"time"

	v1 "github.com/turnforge/weewar/gen/go/weewar/v1/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UsersService defines the interface for user management operations
type UsersService interface {
	// Create a new user profile
	CreateUser(context.Context, *v1.CreateUserRequest) (*v1.CreateUserResponse, error)
	// Batch get multiple users by ID
	GetUsers(context.Context, *v1.GetUsersRequest) (*v1.GetUsersResponse, error)
	// List users with pagination
	ListUsers(context.Context, *v1.ListUsersRequest) (*v1.ListUsersResponse, error)
	// Get a specific user by ID
	GetUser(context.Context, *v1.GetUserRequest) (*v1.GetUserResponse, error)
	// Delete a user
	DeleteUser(context.Context, *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error)
	// Update a user profile
	UpdateUser(context.Context, *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error)
	// EnsureUser creates or updates a user profile (used after auth)
	EnsureUser(ctx context.Context, userId string, name string, email string, imageUrl string) (*v1.User, error)
}

// UserStorageProvider is implemented by concrete backends (fsbe, gormbe, gaebe)
// to provide raw storage operations for users
type UserStorageProvider interface {
	// Read operations
	LoadUser(ctx context.Context, id string) (*v1.User, error)
	ListAllUsers(ctx context.Context) ([]*v1.User, error)

	// Write operations
	SaveUser(ctx context.Context, id string, user *v1.User) error
	DeleteFromStorage(ctx context.Context, id string) error

	// Check if user exists
	UserExists(ctx context.Context, id string) bool
}

// BaseUsersService provides shared logic for user services
type BaseUsersService struct {
	Self            UsersService        // The actual implementation
	StorageProvider UserStorageProvider // Set by concrete implementations

	// Optional in-memory cache
	CacheEnabled bool
	userCache    map[string]*v1.User
	cacheMu      sync.RWMutex
}

// InitializeCache sets up the in-memory cache
func (s *BaseUsersService) InitializeCache() {
	s.CacheEnabled = true
	s.userCache = make(map[string]*v1.User)
}

// GetUser returns a specific user by ID
func (s *BaseUsersService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	if req.Id == "" {
		return nil, ErrUserIDRequired
	}

	// Check cache first
	if s.CacheEnabled {
		s.cacheMu.RLock()
		if user, ok := s.userCache[req.Id]; ok {
			s.cacheMu.RUnlock()
			return &v1.GetUserResponse{User: user}, nil
		}
		s.cacheMu.RUnlock()
	}

	user, err := s.StorageProvider.LoadUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Update cache
	if s.CacheEnabled {
		s.cacheMu.Lock()
		s.userCache[req.Id] = user
		s.cacheMu.Unlock()
	}

	return &v1.GetUserResponse{User: user}, nil
}

// GetUsers returns multiple users by ID
func (s *BaseUsersService) GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersResponse, error) {
	resp := &v1.GetUsersResponse{
		Users: make(map[string]*v1.User),
	}

	for _, id := range req.Ids {
		user, err := s.StorageProvider.LoadUser(ctx, id)
		if err == nil && user != nil {
			resp.Users[id] = user
		}
	}

	return resp, nil
}

// ListUsers returns all users with pagination
func (s *BaseUsersService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	users, err := s.StorageProvider.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.ListUsersResponse{
		Items: users,
		Pagination: &v1.PaginationResponse{
			HasMore:      false,
			TotalResults: int32(len(users)),
		},
	}, nil
}

// DeleteUser removes a user
func (s *BaseUsersService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, ErrUserIDRequired
	}

	err := s.StorageProvider.DeleteFromStorage(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.CacheEnabled {
		s.cacheMu.Lock()
		delete(s.userCache, req.Id)
		s.cacheMu.Unlock()
	}

	return &v1.DeleteUserResponse{}, nil
}

// EnsureUser creates or updates a user profile
func (s *BaseUsersService) EnsureUser(ctx context.Context, userId string, name string, email string, imageUrl string) (*v1.User, error) {
	now := time.Now()

	// Check if user exists
	existingUser, err := s.StorageProvider.LoadUser(ctx, userId)
	if err == nil && existingUser != nil {
		// Update existing user
		if name != "" {
			existingUser.Name = name
		}
		if email != "" {
			existingUser.Email = email
		}
		if imageUrl != "" {
			existingUser.ImageUrl = imageUrl
		}
		existingUser.UpdatedAt = timestamppb.New(now)

		if err := s.StorageProvider.SaveUser(ctx, userId, existingUser); err != nil {
			return nil, err
		}

		// Update cache
		if s.CacheEnabled {
			s.cacheMu.Lock()
			s.userCache[userId] = existingUser
			s.cacheMu.Unlock()
		}

		return existingUser, nil
	}

	// Create new user
	user := &v1.User{
		Id:        userId,
		Name:      name,
		Email:     email,
		ImageUrl:  imageUrl,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	if err := s.StorageProvider.SaveUser(ctx, userId, user); err != nil {
		return nil, err
	}

	// Update cache
	if s.CacheEnabled {
		s.cacheMu.Lock()
		s.userCache[userId] = user
		s.cacheMu.Unlock()
	}

	return user, nil
}

// Error types
var (
	ErrUserIDRequired = &UserError{Message: "user ID is required"}
	ErrUserNotFound   = &UserError{Message: "user not found"}
)

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}
