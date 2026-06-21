//go:build !wasm
// +build !wasm

package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/panyam/oneauth/httpauth"
)

const (
	// DevLoginParam is the query parameter the middleware inspects when the
	// gate is on. Visiting any URL with this param sets the session subject
	// to the param's value, with no password or OAuth round-trip.
	DevLoginParam = "dev_user"

	// DevLoginEnvVar is the env var that gates the middleware. Anything
	// other than "true" disables fake login entirely — the middleware
	// short-circuits on the env check and passes the request through
	// unchanged. Default state is OFF; production deploys must keep it
	// that way.
	DevLoginEnvVar = "ENABLE_DEV_FAKE_LOGIN"
)

// DevLoginEnabled reports whether the gate env var is currently set to
// "true". Read at request time, not at startup, so changing the env in dev
// (e.g. via `make devloop ENABLE_DEV_FAKE_LOGIN=true`) takes effect on the
// next request without restarting tests.
func DevLoginEnabled() bool {
	return os.Getenv(DevLoginEnvVar) == "true"
}

// WrapDevLogin returns a middleware that intercepts ?dev_user=<handle> on
// any path and calls oneauth.SetLoggedInSubject with that handle. The
// handle becomes the session subject, the loggedInSubject cookie value,
// and the JWT cookie's `sub` claim — identical to what a real login would
// produce, just without the password / OAuth round-trip.
//
// Callers decide whether to install this middleware based on
// DevLoginEnabled() at handler-build time, so production handlers never
// even hold a reference to it. Wire it via:
//
//	withDevLogin := routes
//	if server.DevLoginEnabled() { withDevLogin = WrapDevLogin(oneauth)(routes) }
//	Session.LoadAndSave(withDevLogin)
//
// The middleware MUST run inside Session.LoadAndSave so the session writes
// performed by SetLoggedInSubject reach the SCS store.
//
// Every fake-login event emits an slog.Warn line so a production deploy
// that somehow installed this middleware cannot go unnoticed during log
// review.
func WrapDevLogin(oneauth *httpauth.OneAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handle := r.URL.Query().Get(DevLoginParam); handle != "" {
				slog.Warn("dev-mode fake login (NEVER ENABLE IN PROD)",
					"subject", handle,
					"path", r.URL.Path,
					"remote", r.RemoteAddr,
					"env_var", DevLoginEnvVar)
				oneauth.SetLoggedInSubject(handle, w, r)
			}
			next.ServeHTTP(w, r)
		})
	}
}
