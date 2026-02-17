package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
)

// ForgotPasswordPage renders the forgot password form.
type ForgotPasswordPage struct {
	BasePage
	CustomHeader bool
	Sent         bool
}

func (p *ForgotPasswordPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*LilBattleApp]) (err error, finished bool) {
	p.CustomHeader = true
	p.Sent = r.URL.Query().Get("sent") == "true"
	return
}

// ResetPasswordPage renders the password reset form.
type ResetPasswordPage struct {
	BasePage
	CustomHeader bool
	Token        string
	InvalidToken bool
	Success      bool
	Error        string
}

func (p *ResetPasswordPage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*LilBattleApp]) (err error, finished bool) {
	p.CustomHeader = true
	p.Token = r.URL.Query().Get("token")
	p.Success = r.URL.Query().Get("success") == "true"
	p.Error = r.URL.Query().Get("error")
	if p.Token == "" && !p.Success {
		p.InvalidToken = true
	}
	return
}
