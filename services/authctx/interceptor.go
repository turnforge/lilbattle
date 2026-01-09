package authctx

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptorConfig configures the auth interceptor behavior.
type AuthInterceptorConfig struct {
	// RequireAuth when true rejects unauthenticated requests.
	// When false, requests proceed but UserIDFromContext returns empty.
	RequireAuth bool

	// PublicMethods is a list of method names that don't require auth.
	// Only used when RequireAuth is true.
	PublicMethods map[string]bool
}

// DefaultAuthInterceptorConfig returns a config that requires auth for all methods.
func DefaultAuthInterceptorConfig() *AuthInterceptorConfig {
	return &AuthInterceptorConfig{
		RequireAuth:   true,
		PublicMethods: make(map[string]bool),
	}
}

// UnaryAuthInterceptor returns a gRPC unary interceptor that processes auth metadata.
// It handles the X-Switch-User header when ENABLE_SWITCH_AUTH is set.
func UnaryAuthInterceptor(config *AuthInterceptorConfig) grpc.UnaryServerInterceptor {
	if config == nil {
		config = DefaultAuthInterceptorConfig()
	}

	switchAuthEnabled := os.Getenv("ENABLE_SWITCH_AUTH") == "true"

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract user ID from metadata
		userID := ""
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// Check for switch user first (only if enabled)
			if switchAuthEnabled {
				if values := md.Get(MetadataKeySwitchUser); len(values) > 0 && values[0] != "" {
					userID = values[0]
				}
			}

			// Fall back to actual user ID
			if userID == "" {
				if values := md.Get(MetadataKeyUserID); len(values) > 0 {
					userID = values[0]
				}
			}
		}

		// Check if auth is required for this method
		if config.RequireAuth && !config.PublicMethods[info.FullMethod] {
			if userID == "" {
				return nil, status.Error(codes.Unauthenticated, "authentication required")
			}
		}

		return handler(ctx, req)
	}
}

// StreamAuthInterceptor returns a gRPC stream interceptor that processes auth metadata.
func StreamAuthInterceptor(config *AuthInterceptorConfig) grpc.StreamServerInterceptor {
	if config == nil {
		config = DefaultAuthInterceptorConfig()
	}

	switchAuthEnabled := os.Getenv("ENABLE_SWITCH_AUTH") == "true"

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		// Extract user ID from metadata
		userID := ""
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// Check for switch user first (only if enabled)
			if switchAuthEnabled {
				if values := md.Get(MetadataKeySwitchUser); len(values) > 0 && values[0] != "" {
					userID = values[0]
				}
			}

			// Fall back to actual user ID
			if userID == "" {
				if values := md.Get(MetadataKeyUserID); len(values) > 0 {
					userID = values[0]
				}
			}
		}

		// Check if auth is required for this method
		if config.RequireAuth && !config.PublicMethods[info.FullMethod] {
			if userID == "" {
				return status.Error(codes.Unauthenticated, "authentication required")
			}
		}

		return handler(srv, ss)
	}
}

// NewPublicMethodsConfig creates a config with the specified public methods.
func NewPublicMethodsConfig(publicMethods ...string) *AuthInterceptorConfig {
	config := &AuthInterceptorConfig{
		RequireAuth:   true,
		PublicMethods: make(map[string]bool),
	}
	for _, method := range publicMethods {
		config.PublicMethods[method] = true
	}
	return config
}

// OptionalAuthConfig returns a config that allows unauthenticated requests.
func OptionalAuthConfig() *AuthInterceptorConfig {
	return &AuthInterceptorConfig{
		RequireAuth:   false,
		PublicMethods: make(map[string]bool),
	}
}
