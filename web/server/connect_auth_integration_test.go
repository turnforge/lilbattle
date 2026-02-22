//go:build !wasm
// +build !wasm

package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	oa "github.com/panyam/oneauth"
	oagrpc "github.com/panyam/oneauth/grpc"
	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	v1s "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/services"
	v1connect "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/services/lilbattlev1connect"
	"github.com/turnforge/lilbattle/services/connectclient"
	"google.golang.org/grpc"
)

// recordingWorldsClient is a mock WorldsServiceClient that records the user ID
// extracted from the gRPC incoming metadata on each call.
type recordingWorldsClient struct {
	mu      sync.Mutex
	userIDs []string
}

func (r *recordingWorldsClient) recordUserID(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userIDs = append(r.userIDs, oagrpc.UserIDFromContext(ctx))
}

func (r *recordingWorldsClient) lastUserID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.userIDs) == 0 {
		return ""
	}
	return r.userIDs[len(r.userIDs)-1]
}

func (r *recordingWorldsClient) CreateWorld(ctx context.Context, in *v1.CreateWorldRequest, opts ...grpc.CallOption) (*v1.CreateWorldResponse, error) {
	r.recordUserID(ctx)
	return &v1.CreateWorldResponse{}, nil
}

func (r *recordingWorldsClient) GetWorlds(ctx context.Context, in *v1.GetWorldsRequest, opts ...grpc.CallOption) (*v1.GetWorldsResponse, error) {
	r.recordUserID(ctx)
	return &v1.GetWorldsResponse{}, nil
}

func (r *recordingWorldsClient) ListWorlds(ctx context.Context, in *v1.ListWorldsRequest, opts ...grpc.CallOption) (*v1.ListWorldsResponse, error) {
	r.recordUserID(ctx)
	return &v1.ListWorldsResponse{}, nil
}

func (r *recordingWorldsClient) GetWorld(ctx context.Context, in *v1.GetWorldRequest, opts ...grpc.CallOption) (*v1.GetWorldResponse, error) {
	r.recordUserID(ctx)
	return &v1.GetWorldResponse{}, nil
}

func (r *recordingWorldsClient) DeleteWorld(ctx context.Context, in *v1.DeleteWorldRequest, opts ...grpc.CallOption) (*v1.DeleteWorldResponse, error) {
	r.recordUserID(ctx)
	return &v1.DeleteWorldResponse{}, nil
}

func (r *recordingWorldsClient) UpdateWorld(ctx context.Context, in *v1.UpdateWorldRequest, opts ...grpc.CallOption) (*v1.UpdateWorldResponse, error) {
	r.recordUserID(ctx)
	return &v1.UpdateWorldResponse{}, nil
}

// Compile-time check that recordingWorldsClient implements WorldsServiceClient.
var _ v1s.WorldsServiceClient = (*recordingWorldsClient)(nil)

type authTestServer struct {
	server   *httptest.Server
	recorder *recordingWorldsClient
}

// setupAuthTestServer wires:
//
//	recordingWorldsClient → ConnectWorldsServiceAdapter → v1connect.NewWorldsServiceHandler
//	  → oa.Middleware.ExtractUser → http.ServeMux → httptest.NewServer
//
// VerifyToken: tokens matching "valid-token-{userID}" return userID; others error.
func setupAuthTestServer(t *testing.T) *authTestServer {
	t.Helper()

	recorder := &recordingWorldsClient{}
	adapter := NewConnectWorldsServiceAdapter(recorder)
	path, handler := v1connect.NewWorldsServiceHandler(adapter)

	authMiddleware := &oa.Middleware{
		VerifyToken: func(tokenString string) (string, any, error) {
			if userID, ok := strings.CutPrefix(tokenString, "valid-token-"); ok {
				return userID, nil, nil
			}
			return "", nil, fmt.Errorf("invalid token")
		},
	}

	mux := http.NewServeMux()
	mux.Handle(path, authMiddleware.ExtractUser(handler))

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &authTestServer{server: server, recorder: recorder}
}

func TestConnectAuth_CreateWorldWithAuth(t *testing.T) {
	ts := setupAuthTestServer(t)
	client := connectclient.NewConnectWorldsClientWithAuth(ts.server.URL, "valid-token-user-123")

	_, err := client.CreateWorld(context.Background(), &v1.CreateWorldRequest{})
	if err != nil {
		t.Fatalf("CreateWorld failed: %v", err)
	}

	if got := ts.recorder.lastUserID(); got != "user-123" {
		t.Errorf("userID = %q, want %q", got, "user-123")
	}
}

func TestConnectAuth_ListWorldsNoAuth(t *testing.T) {
	ts := setupAuthTestServer(t)
	client := connectclient.NewConnectWorldsClient(ts.server.URL)

	_, err := client.ListWorlds(context.Background(), &v1.ListWorldsRequest{})
	if err != nil {
		t.Fatalf("ListWorlds failed: %v", err)
	}

	if got := ts.recorder.lastUserID(); got != "" {
		t.Errorf("userID = %q, want empty string", got)
	}
}

func TestConnectAuth_InvalidToken(t *testing.T) {
	ts := setupAuthTestServer(t)
	client := connectclient.NewConnectWorldsClientWithAuth(ts.server.URL, "garbage")

	_, err := client.ListWorlds(context.Background(), &v1.ListWorldsRequest{})
	if err != nil {
		t.Fatalf("ListWorlds failed: %v", err)
	}

	if got := ts.recorder.lastUserID(); got != "" {
		t.Errorf("userID = %q, want empty string", got)
	}
}

func TestConnectAuth_MultipleUsers(t *testing.T) {
	ts := setupAuthTestServer(t)

	aliceClient := connectclient.NewConnectWorldsClientWithAuth(ts.server.URL, "valid-token-alice")
	_, err := aliceClient.CreateWorld(context.Background(), &v1.CreateWorldRequest{})
	if err != nil {
		t.Fatalf("CreateWorld (alice) failed: %v", err)
	}
	if got := ts.recorder.lastUserID(); got != "alice" {
		t.Errorf("alice: userID = %q, want %q", got, "alice")
	}

	bobClient := connectclient.NewConnectWorldsClientWithAuth(ts.server.URL, "valid-token-bob")
	_, err = bobClient.CreateWorld(context.Background(), &v1.CreateWorldRequest{})
	if err != nil {
		t.Fatalf("CreateWorld (bob) failed: %v", err)
	}
	if got := ts.recorder.lastUserID(); got != "bob" {
		t.Errorf("bob: userID = %q, want %q", got, "bob")
	}
}

func TestConnectAuth_AllWriteEndpoints(t *testing.T) {
	ts := setupAuthTestServer(t)
	client := connectclient.NewConnectWorldsClientWithAuth(ts.server.URL, "valid-token-user-1")

	tests := []struct {
		name string
		call func() error
	}{
		{"CreateWorld", func() error {
			_, err := client.CreateWorld(context.Background(), &v1.CreateWorldRequest{})
			return err
		}},
		{"UpdateWorld", func() error {
			_, err := client.UpdateWorld(context.Background(), &v1.UpdateWorldRequest{})
			return err
		}},
		{"DeleteWorld", func() error {
			_, err := client.DeleteWorld(context.Background(), &v1.DeleteWorldRequest{})
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
			if got := ts.recorder.lastUserID(); got != "user-1" {
				t.Errorf("%s: userID = %q, want %q", tt.name, got, "user-1")
			}
		})
	}
}
