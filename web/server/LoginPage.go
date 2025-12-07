package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
)

type LoginConfig struct {
	EnableEmailLogin     bool
	EnableGoogleLogin    bool
	EnableGitHubLogin    bool
	EnableMicrosoftLogin bool
	EnableAppleLogin     bool
}

type LoginPage struct {
	BasePage
	Header      Header
	CallbackURL string
	CsrfToken   string
	Config      LoginConfig
}

type RegisterPage struct {
	BasePage
	Header         Header
	CallbackURL    string
	CsrfToken      string
	Name           string
	Email          string
	Password       string
	VerifyPassword string
	Errors         map[string]string
}

func (p *LoginPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	p.DisableSplashScreen = true
	err, finished = p.Header.Load(r, w, app)
	p.CallbackURL = r.URL.Query().Get("callbackURL")

	p.Config = LoginConfig{
		EnableEmailLogin:     true,
		EnableGoogleLogin:    true,
		EnableGitHubLogin:    true,
		EnableMicrosoftLogin: false,
		EnableAppleLogin:     false,
	}
	return
}

func (p *RegisterPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	err, finished = p.Header.Load(r, w, app)
	p.CallbackURL = r.URL.Query().Get("callbackURL")
	return
}
