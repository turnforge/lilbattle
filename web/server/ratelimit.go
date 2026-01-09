package server

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter provides a simple in-memory rate limiting using sliding window.
// For production with multiple instances, consider using Redis-based rate limiting.
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int           // max requests per window
	window   time.Duration // time window
}

// NewRateLimiter creates a new rate limiter with the specified limit per window.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	// Start cleanup goroutine
	go rl.cleanup()
	return rl
}

// Allow checks if a request from the given key should be allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Filter out old requests
	var validRequests []time.Time
	for _, t := range rl.requests[key] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= rl.limit {
		rl.requests[key] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests
	return true
}

// cleanup periodically removes old entries to prevent memory leaks.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for key, times := range rl.requests {
			var validRequests []time.Time
			for _, t := range times {
				if t.After(windowStart) {
					validRequests = append(validRequests, t)
				}
			}
			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitConfig holds configuration for different rate limit tiers.
type RateLimitConfig struct {
	// Auth endpoints (login, signup, forgot-password) - stricter limits
	AuthLimit  int
	AuthWindow time.Duration

	// API endpoints - more permissive
	APILimit  int
	APIWindow time.Duration
}

// DefaultRateLimitConfig returns sensible defaults for rate limiting.
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		AuthLimit:  10,              // 10 auth attempts
		AuthWindow: 15 * time.Minute, // per 15 minutes
		APILimit:   100,             // 100 API calls
		APIWindow:  time.Minute,     // per minute
	}
}

// RateLimitMiddleware creates HTTP middleware that applies rate limiting.
type RateLimitMiddleware struct {
	authLimiter *RateLimiter
	apiLimiter  *RateLimiter
	keyFunc     func(*http.Request) string
}

// NewRateLimitMiddleware creates rate limiting middleware with the given config.
func NewRateLimitMiddleware(config *RateLimitConfig) *RateLimitMiddleware {
	if config == nil {
		config = DefaultRateLimitConfig()
	}
	return &RateLimitMiddleware{
		authLimiter: NewRateLimiter(config.AuthLimit, config.AuthWindow),
		apiLimiter:  NewRateLimiter(config.APILimit, config.APIWindow),
		keyFunc:     defaultKeyFunc,
	}
}

// defaultKeyFunc extracts the client identifier from the request.
// Uses X-Forwarded-For for proxied requests, otherwise RemoteAddr.
func defaultKeyFunc(r *http.Request) string {
	// Check for forwarded IP (common behind load balancers/proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	return r.RemoteAddr
}

// WrapAuth wraps an auth handler with rate limiting.
func (m *RateLimitMiddleware) WrapAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := m.keyFunc(r)
		if !m.authLimiter.Allow(key) {
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// WrapAPI wraps an API handler with rate limiting.
func (m *RateLimitMiddleware) WrapAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := m.keyFunc(r)
		if !m.apiLimiter.Allow(key) {
			http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
