package server

import (
	"context"
	"os"

	"github.com/panyam/oneauth/accounts"
	"github.com/panyam/oneauth/localauth"
	svc "github.com/turnforge/lilbattle/services"
	"golang.org/x/oauth2"
)

// testAuthEnabled returns true if test authentication is enabled via env var.
// This should only be enabled in development/testing environments.
func testAuthEnabled() bool {
	return os.Getenv("ENABLE_TEST_AUTH") == "true"
}

// testUser returns the mock test user for development/testing.
func testUser() *svc.User {
	return &svc.User{
		ID: "test1",
		ProfileInfo: svc.StringMapField{
			Properties: map[string]any{
				"Name": "Test User",
			},
		},
	}
}

func (n *LilBattleApp) GetUserByID(userId string) (accounts.User, error) {
	// Test user bypass - only if ENABLE_TEST_AUTH is set
	if testAuthEnabled() && userId == "test1" {
		return testUser(), nil
	}
	resp, err := n.ClientMgr.GetAuthService().GetUserById(context.Background(), &accounts.GetUserByIDRequest{UserID: userId})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (n *LilBattleApp) EnsureAuthUser(authtype string, provider string, token *oauth2.Token, userInfo map[string]any) (accounts.User, error) {
	// Test user bypass - only if ENABLE_TEST_AUTH is set
	if testAuthEnabled() {
		if email, ok := userInfo["email"].(string); ok && email == "test@gmail.com" {
			return testUser(), nil
		}
	}

	// Assign a random nickname if not already set
	if _, hasNickname := userInfo["nickname"]; !hasNickname {
		userInfo["nickname"] = GenerateRandomNickname()
	}

	return n.ClientMgr.GetAuthService().EnsureAuthUser(authtype, provider, token, userInfo)
}

func (n *LilBattleApp) ValidateUsernamePassword(username string, password string) (accounts.User, error) {
	// Test user bypass - only if ENABLE_TEST_AUTH is set
	if testAuthEnabled() && username == "test@gmail.com" {
		return testUser(), nil
	}
	// For production, delegate to auth service
	usernameType := localauth.DetectUsernameType(username)
	return n.ClientMgr.GetAuthService().ValidateLocalCredentials(username, password, usernameType)
}
