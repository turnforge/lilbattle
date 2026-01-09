package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	oa "github.com/panyam/oneauth"
	"golang.org/x/oauth2"
)

// TwitterOAuth2 implements OAuth 2.0 with PKCE for Twitter/X authentication.
type TwitterOAuth2 struct {
	ClientId       string
	ClientSecret   string
	CallbackURL    string
	HandleUser     oa.HandleUserFunc
	AuthFailureUrl string
	oauthConfig    oauth2.Config
	mux            *http.ServeMux
}

// NewTwitterOAuth2 creates a new Twitter OAuth2 handler.
// Empty parameters will be read from environment variables.
func NewTwitterOAuth2(clientId, clientSecret, callbackUrl string, handleUser oa.HandleUserFunc) *TwitterOAuth2 {
	if clientId == "" {
		clientId = os.Getenv("OAUTH2_TWITTER_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = os.Getenv("OAUTH2_TWITTER_CLIENT_SECRET")
	}
	if callbackUrl == "" {
		callbackUrl = os.Getenv("OAUTH2_TWITTER_CALLBACK_URL")
	}

	t := &TwitterOAuth2{
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		CallbackURL:    callbackUrl,
		HandleUser:     handleUser,
		AuthFailureUrl: "/login?error=twitter_auth_failed",
		oauthConfig: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  callbackUrl,
			Scopes:       []string{"users.read", "tweet.read", "offline.access"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://twitter.com/i/oauth2/authorize",
				TokenURL:  "https://api.x.com/2/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
		mux: http.NewServeMux(),
	}
	t.setupHandlers()
	return t
}

func (t *TwitterOAuth2) Handler() http.Handler {
	return t.mux
}

func (t *TwitterOAuth2) setupHandlers() {
	t.mux.HandleFunc("/", t.handleAuth)
	t.mux.HandleFunc("/callback/", t.handleCallback)
}

// handleAuth initiates the OAuth flow with PKCE.
func (t *TwitterOAuth2) handleAuth(w http.ResponseWriter, r *http.Request) {
	// Generate PKCE verifier
	verifier := oauth2.GenerateVerifier()

	// Store verifier in cookie for callback
	http.SetCookie(w, &http.Cookie{
		Name:     "twitter_oauth_verifier",
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600, // 10 minutes
	})

	// Generate state for CSRF protection
	state := oauth2.GenerateVerifier() // Reuse verifier generation for random state

	// Store state in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "twitter_oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})

	// Store callback URL if provided
	callbackURL := r.URL.Query().Get("callbackURL")
	if callbackURL != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "twitter_callback_url",
			Value:    callbackURL,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})
	}

	// Generate authorization URL with PKCE
	url := t.oauthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleCallback processes the OAuth callback from Twitter.
func (t *TwitterOAuth2) handleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state
	stateCookie, err := r.Cookie("twitter_oauth_state")
	if err != nil {
		log.Printf("Twitter OAuth: missing state cookie: %v", err)
		http.Redirect(w, r, t.AuthFailureUrl, http.StatusTemporaryRedirect)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		log.Printf("Twitter OAuth: state mismatch")
		http.Redirect(w, r, t.AuthFailureUrl, http.StatusTemporaryRedirect)
		return
	}

	// Check for error from Twitter
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		log.Printf("Twitter OAuth error: %s - %s", errMsg, r.URL.Query().Get("error_description"))
		http.Redirect(w, r, t.AuthFailureUrl, http.StatusTemporaryRedirect)
		return
	}

	// Get verifier from cookie
	verifierCookie, err := r.Cookie("twitter_oauth_verifier")
	if err != nil {
		log.Printf("Twitter OAuth: missing verifier cookie: %v", err)
		http.Redirect(w, r, t.AuthFailureUrl, http.StatusTemporaryRedirect)
		return
	}

	// Exchange code for token with PKCE verifier
	code := r.URL.Query().Get("code")
	token, err := t.oauthConfig.Exchange(
		context.Background(),
		code,
		oauth2.VerifierOption(verifierCookie.Value),
	)
	if err != nil {
		log.Printf("Twitter OAuth: token exchange failed: %v", err)
		http.Redirect(w, r, t.AuthFailureUrl, http.StatusTemporaryRedirect)
		return
	}

	// Fetch user data from Twitter
	userData, err := t.getUserData(token.AccessToken)
	if err != nil {
		log.Printf("Twitter OAuth: failed to get user data: %v", err)
		http.Redirect(w, r, t.AuthFailureUrl, http.StatusTemporaryRedirect)
		return
	}

	// Clear OAuth cookies
	for _, name := range []string{"twitter_oauth_verifier", "twitter_oauth_state"} {
		http.SetCookie(w, &http.Cookie{
			Name:   name,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	// Get callback URL and clear cookie
	if callbackCookie, err := r.Cookie("twitter_callback_url"); err == nil {
		// Store callback URL in query params for the handler
		q := r.URL.Query()
		q.Set("callbackURL", callbackCookie.Value)
		r.URL.RawQuery = q.Encode()
		http.SetCookie(w, &http.Cookie{
			Name:   "twitter_callback_url",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	// Call the user handler with the correct signature
	// HandleUserFunc expects: (channelType, channelId string, token *oauth2.Token, profile map[string]any, w ResponseWriter, r *Request)
	channelId := userData["id"].(string)
	t.HandleUser("twitter", channelId, token, userData, w, r)
}

// TwitterUserResponse represents the Twitter API v2 user response.
type TwitterUserResponse struct {
	Data struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Username        string `json:"username"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

// getUserData fetches user profile from Twitter API v2.
func (t *TwitterOAuth2) getUserData(accessToken string) (map[string]any, error) {
	req, err := http.NewRequest("GET", "https://api.x.com/2/users/me?user.fields=profile_image_url", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twitter API error: %s - %s", resp.Status, string(body))
	}

	var twitterUser TwitterUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&twitterUser); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Convert to map for OneAuth compatibility
	userData := map[string]any{
		"id":                twitterUser.Data.ID,
		"name":              twitterUser.Data.Name,
		"username":          twitterUser.Data.Username,
		"profile_image_url": twitterUser.Data.ProfileImageURL,
		// Twitter doesn't provide email in basic scope
		// Use Twitter ID as the identifier
		"email": twitterUser.Data.Username + "@twitter.local",
	}

	return userData, nil
}
