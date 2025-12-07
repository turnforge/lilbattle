package server

import (
	"net/http"

	goal "github.com/panyam/goapplib"
	oa "github.com/panyam/oneauth"
)

type ProfilePage struct {
	BasePage
	Header Header

	// User information
	User              oa.User
	UserID            string
	Email             string
	EmailVerified     bool
	Username          string
	Profile           map[string]any
	VerificationSent  bool
	VerificationError string
}

func (p *ProfilePage) Load(r *http.Request, w http.ResponseWriter, app *goal.App[*WeewarApp]) (err error, finished bool) {
	p.Title = "Profile"
	p.ActiveTab = "profile"
	p.DisableSplashScreen = true

	err, finished = p.Header.Load(r, w, app)
	if err != nil || finished {
		return
	}

	ctx := app.Context
	p.UserID = ctx.AuthMiddleware.GetLoggedInUserId(r)
	if p.UserID == "" {
		http.Redirect(w, r, "/login?callbackURL=/profile", http.StatusFound)
		return nil, true
	}

	if ctx.AuthService != nil {
		p.User, err = ctx.AuthService.GetUserById(p.UserID)
		if err != nil || p.User == nil {
			http.Redirect(w, r, "/login?callbackURL=/profile", http.StatusFound)
			return nil, true
		}

		p.Profile = p.User.Profile()

		if email, ok := p.Profile["email"].(string); ok {
			p.Email = email
		}
		if username, ok := p.Profile["username"].(string); ok {
			p.Username = username
		}

		if p.Email != "" {
			identity, _, identityErr := ctx.AuthService.GetIdentity("email", p.Email, false)
			if identityErr == nil && identity != nil {
				p.EmailVerified = identity.Verified
			}
		}
	}

	if r.URL.Query().Get("verification_sent") == "true" {
		p.VerificationSent = true
	}
	if verifyErr := r.URL.Query().Get("verification_error"); verifyErr != "" {
		p.VerificationError = verifyErr
	}

	return
}
