package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	goalservices "github.com/panyam/goapplib/services"
	"github.com/panyam/oneauth/accounts"
	"github.com/panyam/oneauth/apiauth"
	"github.com/panyam/oneauth/federatedauth"
	"github.com/panyam/oneauth/httpauth"
	"github.com/panyam/oneauth/localauth"
	oa2 "github.com/panyam/oneauth/oauth2"
	oafs "github.com/panyam/oneauth/stores/fs"
)

func newEmailSender() localauth.SendEmail {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("RESEND_API_KEY not set, using console email sender")
		return &localauth.ConsoleEmailSender{}
	}
	fromAddr := os.Getenv("RESEND_FROM_EMAIL")
	if fromAddr == "" {
		fromAddr = "LilBattle <noreply@lilbattle.com>"
	}
	log.Printf("Using Resend email sender (from: %s)", fromAddr)
	return NewResendEmailSender(apiKey, fromAddr)
}

func setupAuthService(session *scs.SessionManager) (*goalservices.AuthService, accounts.UsernameStore, *httpauth.OneAuth) {
	// Initialize authentication
	storagePath := os.Getenv("LILBATTLE_USER_STORAGE_PATH")
	if storagePath == "" {
		storagePath = filepath.Join(os.Getenv("HOME"), "dev-app-data", "lilbattle", "storage")
	}
	authService := goalservices.NewAuthService(storagePath)

	// Create UsernameStore for username → userID mapping
	usernameStore := oafs.NewFSUsernameStore(storagePath)
	authService.UsernameStore = usernameStore

	oneauth := httpauth.New("lilbattle")
	oneauth.Session = session
	oneauth.Middleware.SessionGetter = func(r *http.Request, key string) any {
		return session.GetString(r.Context(), key)
	}

	// OAuthBridge wires the OneAuth session/cookie machinery to AuthService's
	// EnsureAuthUser orchestration. Use bridge.SaveUserAndRedirect as the
	// post-callback handler for each provider and for LocalAuth.
	bridge := federatedauth.NewOAuthBridge(oneauth, authService)

	// OAuth providers - credentials loaded from environment
	oneauth.AddAuth("/google", oa2.NewGoogleOAuth2("", "", "", bridge.SaveUserAndRedirect).Handler())
	oneauth.AddAuth("/github", oa2.NewGithubOAuth2("", "", "", bridge.SaveUserAndRedirect).Handler())
	oneauth.AddAuth("/twitter", NewTwitterOAuth2("", "", "", bridge.SaveUserAndRedirect).Handler())

	// Get base URL for verification/reset links
	baseURL := os.Getenv("LILBATTLE_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Create credentials validator that supports email OR username login
	// - If input contains "@", treats as email
	// - Otherwise, looks up username in UsernameStore to find userID
	validateCredentials := localauth.NewCredentialsValidatorWithUsername(
		authService.IdentityStore,
		authService.ChannelStore,
		authService.UserStore,
		usernameStore,
	)

	// Local authentication (username/password)
	localAuth := &localauth.LocalAuth{
		ValidateCredentials:      validateCredentials,
		CreateUser:               authService.CreateLocalUser,
		ValidateSignup:           nil, // Policy handles validation now
		EmailSender:              newEmailSender(),
		TokenStore:               authService.TokenStore,
		BaseURL:                  baseURL,
		RequireEmailVerification: false,   // Optional verification
		UsernameField:            "email", // Form field name (auto-detection happens after parsing)
		HandleUser:               bridge.SaveUserAndRedirect,
		VerifyEmail:              authService.VerifyEmailByToken,
		UpdatePassword:           authService.UpdatePassword,
		UsernameStore:            usernameStore,

		// Signup policy: email required, username NOT collected at signup
		SignupPolicy: &localauth.SignupPolicy{
			RequireUsername:       false, // Username added later via profile
			RequireEmail:          true,
			RequirePassword:       true,
			EnforceUsernameUnique: false, // Not enforcing at signup since not collected
			EnforceEmailUnique:    true,
			MinPasswordLength:     8,
		},

		// URLs for redirect-based error handling
		LoginURL:  "/login",
		SignupURL: "/login", // Same page, different tab

		// Redirect-based error handling with flash messages
		OnSignupError: func(err *accounts.AuthError, w http.ResponseWriter, r *http.Request) bool {
			session.Put(r.Context(), "auth_error", err.Message)
			session.Put(r.Context(), "auth_error_field", err.Field)
			session.Put(r.Context(), "auth_mode", "signup")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return true
		},
		OnLoginError: func(err *accounts.AuthError, w http.ResponseWriter, r *http.Request) bool {
			session.Put(r.Context(), "auth_error", err.Message)
			session.Put(r.Context(), "auth_error_field", err.Field)
			session.Put(r.Context(), "auth_mode", "login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return true
		},
	}

	oneauth.AddAuth("/login", localAuth)
	oneauth.AddAuth("/signup", http.HandlerFunc(localAuth.HandleSignup))

	// API/CLI token-based authentication. The oneauth v0.1.31 migration
	// replaced the &apiauth.APIAuth{...} god struct with apiauth.NewOneAuth
	// (oneauth #298). Field name shifts: JWTSecretKey → SigningKey + SigningAlg,
	// JWTIssuer → Issuer, JWTAudience → Audience, RefreshTokenStore → RefreshStore,
	// AccessTokenExpiry → AccessExpiry. RefreshTokenExpiry is gone — refresh
	// lifecycle now lives on the RefreshStore impl.
	jwtSecret := os.Getenv("JWT_CLI_SECRET")
	if jwtSecret == "" {
		jwtSecret = "lilbattle-dev-secret-change-in-production" // Dev fallback
	}
	refreshTokenStore := oafs.NewFSRefreshTokenStore(storagePath)
	apiOA := apiauth.NewOneAuth(apiauth.OneAuthConfig{
		SigningKey:          []byte(jwtSecret),
		SigningAlg:          "HS256",
		Issuer:              "lilbattle",
		Audience:            "cli",
		AccessExpiry:        30 * 24 * time.Hour, // 30 days for CLI tokens
		RefreshStore:        refreshTokenStore,
		ValidateCredentials: authService.ValidateLocalCredentials,
	})
	oneauth.AddAuth("/cli/token", apiauth.NewTokenEndpointHandler(apiOA))

	// Wire OneAuth's JWT validation into the HTTP middleware so Bearer tokens
	// from API/CLI clients are accepted by GetLoggedInSubject. Validator is now
	// a struct field on *OneAuth (was a method on *APIAuth).
	oneauth.Middleware.VerifyToken = func(tokenString string) (string, any, error) {
		resp, err := apiOA.Validator.ValidateToken(context.Background(), &apiauth.ValidateTokenRequest{Token: tokenString})
		if err != nil {
			return "", nil, err
		}
		return resp.Info.Subject, resp.Info, nil
	}

	oneauth.AddAuth("/verify-email", http.HandlerFunc(localAuth.HandleVerifyEmail))
	oneauth.AddAuth("/forgot-password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "/forgot-password", http.StatusFound)
			return
		}
		// POST: create token and send email, then redirect
		email := r.FormValue("email")
		if email == "" {
			http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
			return
		}
		createResp, err := authService.TokenStore.CreateToken(r.Context(), &localauth.CreateVerificationTokenRequest{
			Email:          email,
			Type:           localauth.VerificationTypePasswordReset,
			ExpiryDuration: localauth.VerificationExpiryPasswordReset,
		})
		if err != nil {
			log.Printf("Error creating reset token: %v", err)
		} else {
			resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", baseURL, createResp.Token.Token)
			if err := localAuth.EmailSender.SendPasswordResetEmail(email, resetLink); err != nil {
				log.Printf("Error sending reset email: %v", err)
			}
		}
		// Always redirect with sent=true (don't reveal if email exists)
		http.Redirect(w, r, "/forgot-password?sent=true", http.StatusSeeOther)
	}))
	oneauth.AddAuth("/reset-password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			target := "/reset-password"
			if token := r.URL.Query().Get("token"); token != "" {
				target += "?token=" + token
			}
			http.Redirect(w, r, target, http.StatusFound)
			return
		}
		// POST: validate token and update password, then redirect
		token := r.FormValue("token")
		password := r.FormValue("password")
		if token == "" || password == "" {
			http.Redirect(w, r, "/reset-password?error=Token+and+password+required&token="+token, http.StatusSeeOther)
			return
		}
		getResp, err := authService.TokenStore.GetToken(r.Context(), &localauth.GetVerificationTokenRequest{Token: token})
		if err != nil || getResp.Token == nil || getResp.Token.Type != localauth.VerificationTypePasswordReset {
			http.Redirect(w, r, "/reset-password?error=Invalid+or+expired+reset+link", http.StatusSeeOther)
			return
		}
		if len(password) < 8 {
			http.Redirect(w, r, "/reset-password?error=Password+must+be+at+least+8+characters&token="+token, http.StatusSeeOther)
			return
		}
		if err := authService.UpdatePassword(getResp.Token.Email, password); err != nil {
			log.Printf("Error resetting password: %v", err)
			http.Redirect(w, r, "/reset-password?error=Failed+to+reset+password&token="+token, http.StatusSeeOther)
			return
		}
		_, _ = authService.TokenStore.DeleteToken(r.Context(), &localauth.DeleteVerificationTokenRequest{Token: token})
		http.Redirect(w, r, "/reset-password?success=true", http.StatusSeeOther)
	}))

	// Resend verification email
	oneauth.AddAuth("/resend-verification", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		if email == "" {
			http.Redirect(w, r, "/profile?verification_error=Email is required", http.StatusFound)
			return
		}

		// Get the identity to find the user ID
		identityResp, err := authService.IdentityStore.GetIdentity(r.Context(), &accounts.GetIdentityRequest{
			IdentityType:  "email",
			IdentityValue: email,
		})
		if err != nil || identityResp.Identity == nil {
			// For security, don't reveal if email exists - just say success
			http.Redirect(w, r, "/profile?verification_sent=true", http.StatusFound)
			return
		}

		// Create verification token
		createResp, err := authService.TokenStore.CreateToken(r.Context(), &localauth.CreateVerificationTokenRequest{
			Subject:        identityResp.Identity.UserID,
			Email:          email,
			Type:           localauth.VerificationTypeEmail,
			ExpiryDuration: localauth.VerificationExpiryEmail,
		})
		if err != nil {
			log.Printf("Error creating verification token: %v", err)
			http.Redirect(w, r, "/profile?verification_error=Failed to create verification token", http.StatusFound)
			return
		}

		// Send verification email
		verificationLink := baseURL + "/auth/verify-email?token=" + createResp.Token.Token
		if err := localAuth.EmailSender.SendVerificationEmail(email, verificationLink); err != nil {
			log.Printf("Error sending verification email: %v", err)
			http.Redirect(w, r, "/profile?verification_error=Failed to send verification email", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/profile?verification_sent=true", http.StatusFound)
	}))

	// Change password endpoint (for users who already have a password)
	oneauth.AddAuth("/change-password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}

		userId := oneauth.Middleware.GetLoggedInSubject(r)
		if userId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Not logged in"}`))
			return
		}

		userResp, err := authService.UserStore.GetUserById(r.Context(), &accounts.GetUserByIDRequest{UserID: userId})
		if err != nil || userResp.User == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "User not found"}`))
			return
		}

		profile := userResp.User.Profile()
		email, ok := profile["email"].(string)
		if !ok || email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "No email associated with account"}`))
			return
		}

		if err := r.ParseForm(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Invalid form data"}`))
			return
		}

		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")

		if currentPassword == "" || newPassword == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Current password and new password are required"}`))
			return
		}

		// Verify current password
		_, err = authService.ValidateLocalCredentials(email, currentPassword, "email")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Current password is incorrect"}`))
			return
		}

		// Update password
		if err := authService.UpdatePassword(email, newPassword); err != nil {
			log.Printf("Error updating password: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to update password"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))

	// Set password endpoint (for OAuth-only users setting password for the first time)
	oneauth.AddAuth("/set-password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}

		userId := oneauth.Middleware.GetLoggedInSubject(r)
		if userId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Not logged in"}`))
			return
		}

		userResp, err := authService.UserStore.GetUserById(r.Context(), &accounts.GetUserByIDRequest{UserID: userId})
		if err != nil || userResp.User == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "User not found"}`))
			return
		}

		profile := userResp.User.Profile()
		email, ok := profile["email"].(string)
		if !ok || email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "No email associated with account"}`))
			return
		}

		// Check if user already has a password - should use change-password instead
		identityKey := accounts.IdentityKey("email", email)
		channelResp, _ := authService.ChannelStore.GetChannel(r.Context(), &accounts.GetChannelRequest{
			Provider:    "local",
			IdentityKey: identityKey,
		})
		if channelResp != nil && channelResp.Channel != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Password already set. Use change-password endpoint instead."}`))
			return
		}

		if err := r.ParseForm(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Invalid form data"}`))
			return
		}

		newPassword := r.FormValue("new_password")
		if newPassword == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Password is required"}`))
			return
		}

		// Create local channel with password
		config := localauth.LinkLocalCredentialsConfig{
			UserStore:     authService.UserStore,
			IdentityStore: authService.IdentityStore,
			ChannelStore:  authService.ChannelStore,
			UsernameStore: usernameStore,
		}

		// Get username from profile if set
		username, _ := profile["username"].(string)

		if err := localauth.LinkLocalCredentials(config, userId, username, newPassword, email); err != nil {
			log.Printf("Error setting password for user %s: %v", userId, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to set password"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))

	return authService, usernameStore, oneauth
}
