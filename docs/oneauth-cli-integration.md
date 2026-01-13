# OneAuth CLI Integration Requirements

This document specifies the integration points needed in oneauth to support CLI authentication for lilbattle.

## Overview

The CLI needs to authenticate users to remote servers and store credentials locally. This requires oneauth to provide:

1. A CLI token endpoint for username/password authentication
2. JWT tokens with configurable expiration
3. Token validation middleware that accepts Bearer tokens

## Required Endpoints

### 1. CLI Token Endpoint

**Endpoint**: `POST /auth/cli/token`

**Purpose**: Exchange email/password credentials for a JWT token suitable for CLI use.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

**Success Response** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-02-01T00:00:00Z",
  "user_id": "user123",
  "user_email": "user@example.com"
}
```

**Error Responses**:
- 401 Unauthorized: Invalid credentials
- 429 Too Many Requests: Rate limited

**Implementation Notes**:
- Should use the same `ValidateCredentials` function as LocalAuth
- Token expiration should be configurable (default: 30 days for CLI tokens)
- Should be rate-limited (stricter than normal auth endpoints)

### 2. Token Refresh Endpoint (Optional)

**Endpoint**: `POST /auth/cli/refresh`

**Purpose**: Refresh an expiring token without re-entering credentials.

**Request**:
```json
{
  "token": "existing-jwt-token"
}
```

**Success Response** (200 OK):
```json
{
  "token": "new-jwt-token",
  "expires_at": "2025-03-01T00:00:00Z"
}
```

## JWT Token Requirements

### Token Claims

The JWT token should include:

```json
{
  "sub": "user123",           // User ID
  "email": "user@example.com", // User email
  "iat": 1704067200,          // Issued at
  "exp": 1706745600,          // Expiration
  "iss": "lilbattle",         // Issuer
  "aud": "cli",               // Audience (distinguishes CLI tokens)
  "jti": "unique-token-id"    // Token ID for revocation
}
```

### Token Signing

- Use HS256 with a server-side secret, or RS256 with key pairs
- Secret should be configurable via environment variable: `JWT_CLI_SECRET`

## Middleware Integration

### Bearer Token Validation

The existing oneauth middleware should be extended to accept:

```
Authorization: Bearer <jwt-token>
```

**Pseudocode**:
```go
func (m *Middleware) GetLoggedInUserId(r *http.Request) string {
    // First, check session (existing behavior)
    if userId := m.getSessionUserId(r); userId != "" {
        return userId
    }

    // Then, check Bearer token
    authHeader := r.Header.Get("Authorization")
    if strings.HasPrefix(authHeader, "Bearer ") {
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if claims, err := m.validateJWT(token); err == nil {
            return claims.Subject // user ID
        }
    }

    return ""
}
```

### gRPC Metadata

For Connect RPC calls, the middleware should also check:

```go
// From gRPC metadata
md, ok := metadata.FromIncomingContext(ctx)
if ok {
    if tokens := md.Get("authorization"); len(tokens) > 0 {
        // Validate Bearer token from metadata
    }
}
```

## Go Interface Proposal

```go
package oneauth

// CLIAuthConfig configures CLI authentication
type CLIAuthConfig struct {
    // TokenExpiration is how long CLI tokens are valid (default: 30 days)
    TokenExpiration time.Duration

    // JWTSecret is the secret used to sign tokens
    JWTSecret string

    // ValidateCredentials validates email/password
    ValidateCredentials func(email, password string) (userID string, err error)
}

// CLIAuth handles CLI authentication
type CLIAuth struct {
    config CLIAuthConfig
}

// NewCLIAuth creates a new CLI auth handler
func NewCLIAuth(config CLIAuthConfig) *CLIAuth

// Handler returns an http.Handler for /auth/cli/* routes
func (c *CLIAuth) Handler() http.Handler

// CLITokenResponse is returned on successful authentication
type CLITokenResponse struct {
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    UserID    string    `json:"user_id"`
    UserEmail string    `json:"user_email"`
}
```

## Integration in LilBattle

In `web/server/auth.go`:

```go
func setupAuthService(session *scs.SessionManager) (*goalservices.AuthService, *oa.OneAuth) {
    // ... existing setup ...

    // Add CLI auth
    cliAuth := oa.NewCLIAuth(oa.CLIAuthConfig{
        TokenExpiration:     30 * 24 * time.Hour, // 30 days
        JWTSecret:           os.Getenv("JWT_CLI_SECRET"),
        ValidateCredentials: authService.ValidateLocalCredentials,
    })
    oneauth.AddAuth("/cli", cliAuth)

    // ... rest of setup ...
}
```

## Security Considerations

1. **Rate Limiting**: CLI token endpoint should be strictly rate-limited (e.g., 5 attempts per 15 minutes per IP)

2. **Token Storage**: CLI should store tokens with restricted file permissions (0600)

3. **Token Revocation**: Consider adding a token revocation endpoint for security-conscious users

4. **Audit Logging**: Log CLI token issuance for security auditing

5. **Separate Secrets**: Use a different secret for CLI tokens vs session tokens to allow independent rotation
