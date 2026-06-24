//go:build !wasm
// +build !wasm

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/panyam/oneauth/accounts"
	"github.com/panyam/oneauth/apiauth"
	"github.com/panyam/oneauth/core"
	oagrpc "github.com/panyam/oneauth/grpc"
	"github.com/panyam/oneauth/httpauth"
	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
	v1s "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/services"
	v1connect "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/services/lilbattlev1connect"
	"github.com/turnforge/lilbattle/services/connectclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// assertUnauthenticated confirms an error is the auth-rejection class. The
// chain in production looks like:
//
//	oagrpc.UnaryAuthInterceptor returns status.Error(codes.Unauthenticated, ...)
//	  -> connect adapter passes through raw
//	  -> connect client wraps as CodeUnknown but preserves the gRPC message text
//
// So neither status.FromError nor connect.CodeOf alone surfaces
// Unauthenticated for the client-side observer. The durable signal is the
// gRPC status code embedded in the message ("code = Unauthenticated"). Match
// on that. Brittle to a oneauth message change, but the alternative would be
// catching gRPC status pre-Connect wrap, which doesn't exist here without
// custom error mapping in the adapter.
func assertUnauthenticated(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected Unauthenticated, got nil")
	}
	if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
		return
	}
	if connect.CodeOf(err) == connect.CodeUnauthenticated {
		return
	}
	if strings.Contains(err.Error(), "code = Unauthenticated") {
		return
	}
	t.Errorf("expected Unauthenticated; got %v", err)
}

// recordingGamesServer is a stub GamesServiceServer that records the subject
// it observes via incoming gRPC metadata on the auth-required RPC under test.
// It returns canned responses for everything else; the only point is to
// confirm that (a) the call reached the server and (b) the subject arrived
// on the gRPC side, not just in the HTTP context.
type recordingGamesServer struct {
	v1s.UnimplementedGamesServiceServer
	mu         sync.Mutex
	subjects   []string
	callCount  int
}

func (r *recordingGamesServer) GetOptionsAt(ctx context.Context, req *v1.GetOptionsAtRequest) (*v1.GetOptionsAtResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.callCount++
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vs := md.Get(oagrpc.DefaultMetadataKeySubject); len(vs) > 0 {
			r.subjects = append(r.subjects, vs[0])
		} else {
			r.subjects = append(r.subjects, "")
		}
	} else {
		r.subjects = append(r.subjects, "")
	}
	return &v1.GetOptionsAtResponse{}, nil
}

func (r *recordingGamesServer) lastSubject() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.subjects) == 0 {
		return ""
	}
	return r.subjects[len(r.subjects)-1]
}

// authIntegrationFixture wires the full CLI auth handshake:
//
//	connect client (Bearer header)
//	  -> httptest server (apiauth.OneAuth on /cli/token, ConnectGamesService
//	     adapter on /api/...) with httpauth.Middleware extracting the subject
//	  -> grpc client (ClientMgr's grpc.NewClient over loopback)
//	  -> bufconn-bound grpc.Server with oagrpc.UnaryAuthInterceptor in front
//	     of a recordingGamesServer
//
// The point of the fixture is to exercise the boundary the connect.go bug
// crossed: HTTP auth middleware sets subject in core context, the Connect
// adapter forwards it as outgoing gRPC metadata, the real grpc client
// serializes it onto the wire, and the server interceptor reads it from
// incoming metadata. Anything narrower than this misses the bug class.
type authIntegrationFixture struct {
	httpServer *httptest.Server
	grpcServer *grpc.Server
	recorder   *recordingGamesServer
	oa         *apiauth.OneAuth // for minting valid tokens in the test
}

const testJWTSecret = "test-secret-for-cli-auth-integration"

func newAuthIntegrationFixture(t *testing.T) *authIntegrationFixture {
	t.Helper()

	// Boot a real grpc.Server on an OS-assigned port. Use the same interceptor
	// production uses so a metadata-key mismatch (the previous bug class) would
	// fire here.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	recorder := &recordingGamesServer{}
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			oagrpc.UnaryAuthInterceptor(oagrpc.DefaultInterceptorConfig()),
		),
	)
	v1s.RegisterGamesServiceServer(grpcSrv, recorder)
	go func() { _ = grpcSrv.Serve(lis) }()
	t.Cleanup(grpcSrv.Stop)

	// Real grpc client pointed at the test gRPC server. Mirrors
	// ClientMgr.GetGamesSvcClient — production uses grpc.NewClient over a
	// loopback address, so this is the same code path.
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc client: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	gamesClient := v1s.NewGamesServiceClient(conn)

	// Real apiauth.OneAuth for token mint + validation, matching what
	// web/server/auth.go wires in production. A nil RefreshStore is fine
	// for these tests — we exercise access tokens only.
	oa := apiauth.NewOneAuth(apiauth.OneAuthConfig{
		SigningKey:   []byte(testJWTSecret),
		SigningAlg:   "HS256",
		Issuer:       "lilbattle-test",
		Audience:     "cli-test",
		AccessExpiry: 5 * time.Minute,
		ValidateCredentials: func(username, password, usernameType string) (accounts.User, error) {
			if username == "alice" && password == "correct-horse" {
				return &accounts.BasicUser{ID: "alice"}, nil
			}
			return nil, fmt.Errorf("invalid credentials")
		},
	})

	// Auth middleware — production wires VerifyToken through OneAuth's
	// Validator, so do the same here. SubjectGetter from session not needed
	// for CLI flow (Bearer-only).
	authMW := &httpauth.Middleware{
		VerifyToken: func(tokenString string) (string, any, error) {
			resp, err := oa.Validator.ValidateToken(context.Background(), &apiauth.ValidateTokenRequest{Token: tokenString})
			if err != nil {
				return "", nil, err
			}
			return resp.Info.Subject, resp.Info, nil
		},
	}

	// Mirror api.go's wrapWithAuth — same closure shape, same SetSubjectInContext
	// call. If this gets refactored in api.go and not here, the test still
	// passes; the production-equivalence test is in connect_auth_integration_test.go
	// which exercises the same adapter directly. This fixture's job is the
	// HTTP-to-gRPC boundary crossing specifically.
	wrapWithAuth := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subject := authMW.GetLoggedInSubject(r)
			if subject != "" {
				r = r.WithContext(core.SetSubjectInContext(r.Context(), subject))
			}
			handler.ServeHTTP(w, r)
		})
	}

	// Connect adapter wrapping the real grpc client — same wiring as
	// ApiHandler.setupConnectHandlers in production.
	adapter := NewConnectGamesServiceAdapter(gamesClient)
	connectPath, connectHandler := v1connect.NewGamesServiceHandler(adapter)

	mux := http.NewServeMux()
	mux.Handle("/cli/token", apiauth.NewTokenEndpointHandler(oa))
	mux.Handle(connectPath, wrapWithAuth(connectHandler))

	httpServer := httptest.NewServer(mux)
	t.Cleanup(httpServer.Close)

	return &authIntegrationFixture{
		httpServer: httpServer,
		grpcServer: grpcSrv,
		recorder:   recorder,
		oa:         oa,
	}
}

// mintTokenViaCLI POSTs to /cli/token with the password grant (RFC 6749 §4.3)
// and returns the issued access token. The wire shape matches what
// `ww login` does in production — exercising token mint here means a
// regression in apiauth's grant handling would fail this test, not just
// a custom-shape stub.
func (f *authIntegrationFixture) mintTokenViaCLI(t *testing.T, username, password string) string {
	t.Helper()
	form := strings.NewReader("grant_type=password&username=" + username + "&password=" + password)
	req, err := http.NewRequest(http.MethodPost, f.httpServer.URL+"/cli/token", form)
	if err != nil {
		t.Fatalf("build token request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("token endpoint returned %d, want 200", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	if body.AccessToken == "" {
		t.Fatalf("token endpoint returned empty access_token")
	}
	return body.AccessToken
}

// TestCLIAuthIntegration_AuthedRPCSucceeds is the red-before-green case for
// the connect.go outgoing-metadata fix: revert injectAuthMetadata to use
// NewIncomingContext only and this test fails because the gRPC client drops
// the subject at the wire boundary, the server interceptor 401s, and we
// never see "alice" on the recorder.
func TestCLIAuthIntegration_AuthedRPCSucceeds(t *testing.T) {
	f := newAuthIntegrationFixture(t)
	token := f.mintTokenViaCLI(t, "alice", "correct-horse")

	client := connectclient.NewConnectGamesClientWithAuth(f.httpServer.URL, token)
	_, err := client.GetOptionsAt(context.Background(), &v1.GetOptionsAtRequest{})
	if err != nil {
		t.Fatalf("GetOptionsAt with valid Bearer: unexpected error %v", err)
	}
	if got := f.recorder.lastSubject(); got != "alice" {
		t.Errorf("subject at gRPC server = %q, want %q", got, "alice")
	}
}

// TestCLIAuthIntegration_NoBearerIsRejected pins the negative gate. The HTTP
// middleware finds no token, the Connect adapter sends no subject metadata,
// the gRPC interceptor 401s. Recorder must never be called.
func TestCLIAuthIntegration_NoBearerIsRejected(t *testing.T) {
	f := newAuthIntegrationFixture(t)
	client := connectclient.NewConnectGamesClient(f.httpServer.URL) // no auth

	_, err := client.GetOptionsAt(context.Background(), &v1.GetOptionsAtRequest{})
	assertUnauthenticated(t, err)
	if got := f.recorder.callCount; got != 0 {
		t.Errorf("recorder.callCount = %d, want 0 (call should never reach the server)", got)
	}
}

// TestCLIAuthIntegration_WrongSecretIsRejected pins signature verification.
// A token signed with a different secret must not validate, even though it
// is structurally a valid JWT with the right issuer and audience.
func TestCLIAuthIntegration_WrongSecretIsRejected(t *testing.T) {
	f := newAuthIntegrationFixture(t)

	// Independent OneAuth that mints structurally-valid JWTs signed with a
	// secret the production validator does not know.
	attacker := apiauth.NewOneAuth(apiauth.OneAuthConfig{
		SigningKey:   []byte("attacker-secret-not-known-to-server"),
		SigningAlg:   "HS256",
		Issuer:       "lilbattle-test",
		Audience:     "cli-test",
		AccessExpiry: 5 * time.Minute,
	})
	tokResp, err := attacker.Issuer.CreateAccessToken(context.Background(), &apiauth.CreateAccessTokenRequest{
		Subject: "mallory",
	})
	if err != nil {
		t.Fatalf("forge token: %v", err)
	}

	client := connectclient.NewConnectGamesClientWithAuth(f.httpServer.URL, tokResp.Token)
	_, err = client.GetOptionsAt(context.Background(), &v1.GetOptionsAtRequest{})
	assertUnauthenticated(t, err)
	if got := f.recorder.callCount; got != 0 {
		t.Errorf("recorder.callCount = %d, want 0", got)
	}
}
