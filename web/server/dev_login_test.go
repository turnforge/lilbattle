//go:build !wasm
// +build !wasm

package server

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/panyam/oneauth/httpauth"
)

// newDevLoginTestServer builds the smallest server that exercises
// WrapDevLogin end-to-end: an scs session, an oneauth.OneAuth with the
// session wired in (so SetLoggedInSubject can persist), the dev-login
// middleware wrapping a tiny handler that reports the current subject.
// Cookies set by SetLoggedInSubject propagate back via the standard
// httptest.NewServer + http.Client jar pattern.
func newDevLoginTestServer(t *testing.T) (*httptest.Server, *scs.SessionManager, *httpauth.OneAuth) {
	t.Helper()

	session := scs.New()
	oneauth := httpauth.New("test")
	oneauth.Session = session
	oneauth.Middleware.SessionGetter = func(r *http.Request, key string) any {
		return session.GetString(r.Context(), key)
	}

	// The inner handler echoes the current subject so tests can confirm
	// whether the request was authenticated by the middleware.
	echo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subj := oneauth.Middleware.GetLoggedInSubject(r)
		w.Header().Set("X-Test-Subject", subj)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := WrapDevLogin(oneauth)(echo)
	srv := httptest.NewServer(session.LoadAndSave(wrapped))
	t.Cleanup(srv.Close)
	return srv, session, oneauth
}

func newCookieClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookie jar: %v", err)
	}
	return &http.Client{Jar: jar}
}

// TestDevLoginEnabled_RespectsEnvVar pins the gate's read-time semantics:
// the function reflects whatever the env var is set to right now, so
// callers can toggle the gate per-test via t.Setenv without restarting.
func TestDevLoginEnabled_RespectsEnvVar(t *testing.T) {
	t.Setenv(DevLoginEnvVar, "true")
	if !DevLoginEnabled() {
		t.Errorf("expected DevLoginEnabled()=true with env var set; got false")
	}
	t.Setenv(DevLoginEnvVar, "false")
	if DevLoginEnabled() {
		t.Errorf("expected DevLoginEnabled()=false when env var is not exactly 'true'; got true")
	}
	t.Setenv(DevLoginEnvVar, "")
	if DevLoginEnabled() {
		t.Errorf("expected DevLoginEnabled()=false when env var is unset; got true")
	}
}

// TestWrapDevLogin_NoParam pins the no-op path: when the request carries no
// ?dev_user=, the wrapper must pass through without touching the session.
func TestWrapDevLogin_NoParam(t *testing.T) {
	srv, _, _ := newDevLoginTestServer(t)
	client := newCookieClient(t)

	resp, err := client.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	resp.Body.Close()

	if got := resp.Header.Get("X-Test-Subject"); got != "" {
		t.Errorf("expected empty subject without ?dev_user=; got %q", got)
	}
	if hasLoggedInCookie(resp.Cookies()) {
		t.Errorf("expected no loggedInSubject cookie when no ?dev_user= sent")
	}
}

// TestWrapDevLogin_WithParam pins the fake-login path: the handle in
// ?dev_user= becomes the session subject and a loggedInSubject cookie is
// set on the response.
func TestWrapDevLogin_WithParam(t *testing.T) {
	srv, _, _ := newDevLoginTestServer(t)
	client := newCookieClient(t)

	resp, err := client.Get(srv.URL + "/?dev_user=alice")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	resp.Body.Close()

	if got := resp.Header.Get("X-Test-Subject"); got != "alice" {
		t.Errorf("expected subject=alice on the fake-login request itself; got %q", got)
	}
	if !hasLoggedInCookie(resp.Cookies()) {
		t.Errorf("expected loggedInSubject cookie on the response; got %v", resp.Cookies())
	}
}

// TestWrapDevLogin_SubjectPersistsAcrossRequests pins the durable-session
// property — a second request to a URL *without* the param still sees the
// subject set by the first, via the cookie. This is the multi-window
// workflow the issue calls out: visit once with ?dev_user=alice, then
// navigate freely as alice.
func TestWrapDevLogin_SubjectPersistsAcrossRequests(t *testing.T) {
	srv, _, _ := newDevLoginTestServer(t)
	client := newCookieClient(t)

	// First request sets the subject.
	resp1, err := client.Get(srv.URL + "/?dev_user=alice")
	if err != nil {
		t.Fatalf("GET 1: %v", err)
	}
	resp1.Body.Close()

	// Second request, no param — cookie alone must keep the subject.
	resp2, err := client.Get(srv.URL + "/some/other/path")
	if err != nil {
		t.Fatalf("GET 2: %v", err)
	}
	resp2.Body.Close()

	if got := resp2.Header.Get("X-Test-Subject"); got != "alice" {
		t.Errorf("expected subject=alice on the no-param follow-up request; got %q", got)
	}
}

// TestWrapDevLogin_SwapSubject pins the swap workflow — visiting with a
// new ?dev_user= value replaces the previous identity in the session.
func TestWrapDevLogin_SwapSubject(t *testing.T) {
	srv, _, _ := newDevLoginTestServer(t)
	client := newCookieClient(t)

	resp1, _ := client.Get(srv.URL + "/?dev_user=alice")
	resp1.Body.Close()

	resp2, err := client.Get(srv.URL + "/?dev_user=bob")
	if err != nil {
		t.Fatalf("GET swap: %v", err)
	}
	resp2.Body.Close()

	if got := resp2.Header.Get("X-Test-Subject"); got != "bob" {
		t.Errorf("expected subject=bob after swap; got %q", got)
	}
}

func hasLoggedInCookie(cookies []*http.Cookie) bool {
	for _, c := range cookies {
		if c.Name == "loggedInSubject" && c.Value != "" {
			return true
		}
	}
	return false
}
