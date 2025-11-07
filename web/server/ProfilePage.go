package server

import (
	"net/http"

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

func (p *ProfilePage) Load(r *http.Request, w http.ResponseWriter, vc *ViewContext) (err error, finished bool) {
	p.Title = "Profile"
	p.ActiveTab = "profile"
	p.DisableSplashScreen = true

	// Load header
	err, finished = p.Header.Load(r, w, vc)
	if err != nil || finished {
		return
	}

	// Check if user is logged in (request-scoped)
	p.UserID = vc.AuthMiddleware.GetLoggedInUserId(r)
	if p.UserID == "" {
		// Redirect to login if not logged in
		http.Redirect(w, r, "/login?callbackURL=/profile", http.StatusFound)
		return nil, true
	}

	// Load user data from storage (global store)
	if vc.AuthService != nil {
		p.User, err = vc.AuthService.GetUserById(p.UserID)
		if err != nil || p.User == nil {
			// User ID in session but user doesn't exist in storage - redirect to login
			http.Redirect(w, r, "/login?callbackURL=/profile", http.StatusFound)
			return nil, true
		}

		p.Profile = p.User.Profile()

		// Extract email from profile
		if email, ok := p.Profile["email"].(string); ok {
			p.Email = email
		}

		// Extract username
		if username, ok := p.Profile["username"].(string); ok {
			p.Username = username
		}

		// Get email verification status from IdentityStore (global store)
		if p.Email != "" {
			identity, _, identityErr := vc.AuthService.GetIdentity("email", p.Email, false)
			if identityErr == nil && identity != nil {
				p.EmailVerified = identity.Verified
			}
		}
	}

	// Check if we just sent a verification email
	if r.URL.Query().Get("verification_sent") == "true" {
		p.VerificationSent = true
	}
	if verifyErr := r.URL.Query().Get("verification_error"); verifyErr != "" {
		p.VerificationError = verifyErr
	}

	return
}

func (p *ProfilePage) Copy() View {
	return &ProfilePage{}
}
