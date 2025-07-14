package server

import (
	svc "github.com/panyam/turnengine/games/weewar/services"
	oa "github.com/panyam/oneauth"
	"golang.org/x/oauth2"
)

type AuthUser struct {
	*svc.User
}

func (n *AuthUser) Id() string {
	return n.User.Id
}

func (n *AuthUser) Profile() map[string]any {
	return n.User.Profile.Properties
}

func (n *App) GetUserByID(userId string) (oa.User, error) {
	var user AuthUser
	var err error
	if userId == "test1" {
		// Mocking user login
		return &AuthUser{
			User: &svc.User{
				Id: "test1",
				Profile: svc.StringMapField{
					Properties: map[string]any{
						"Name": "Test User",
					},
				},
			},
		}, nil
	}
	user.User, err = n.ClientMgr.GetAuthService().GetUserByID(userId)
	return &user, err
}

func (n *App) EnsureAuthUser(authtype string, provider string, token *oauth2.Token, userInfo map[string]any) (oa.User, error) {
	var user AuthUser
	var err error
	// Mocking user login
	email := userInfo["email"].(string)
	if email == "test@gmail.com" {
		return &AuthUser{
			User: &svc.User{
				Id: "test1",
				Profile: svc.StringMapField{
					Properties: map[string]any{
						"Name": "Test User",
					},
				},
			},
		}, nil
	}
	user.User, err = n.ClientMgr.GetAuthService().EnsureAuthUser(authtype, provider, token, userInfo)
	return &user, err
}

func (n *App) ValidateUsernamePassword(username string, password string) (out oa.User, err error) {
	if username == "test@gmail.com" {
		out = &AuthUser{
			User: &svc.User{
				Id: "test1",
				Profile: svc.StringMapField{
					Properties: map[string]any{
						"Name": "Test User",
					},
				},
			},
		}
	}
	return
}
