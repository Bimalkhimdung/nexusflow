package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is a custom type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey ContextKey = "user_email"
)

// JWTConfig holds JWT middleware configuration
type JWTConfig struct {
	Secret     string
	SkipPaths  []string // Paths that don't require authentication
	ErrorFunc  func(w http.ResponseWriter, r *http.Request, err error)
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RequireAuth creates middleware that requires valid JWT authentication
func RequireAuth(config *JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if path should skip authentication
			if shouldSkipAuth(r.URL.Path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			token, err := extractToken(r)
			if err != nil {
				handleAuthError(w, r, err, config.ErrorFunc)
				return
			}

			// Validate and parse token
			claims, err := validateToken(token, config.Secret)
			if err != nil {
				handleAuthError(w, r, err, config.ErrorFunc)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth creates middleware that optionally extracts user info if token is present
func OptionalAuth(config *JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to extract token
			token, err := extractToken(r)
			if err != nil {
				// No token or invalid format, continue without auth
				next.ServeHTTP(w, r)
				return
			}

			// Try to validate token
			claims, err := validateToken(token, config.Secret)
			if err != nil {
				// Invalid token, continue without auth
				next.ServeHTTP(w, r)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken extracts JWT token from Authorization header
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", &AuthError{Message: "missing authorization header"}
	}

	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", &AuthError{Message: "invalid authorization header format"}
	}

	return parts[1], nil
}

// validateToken validates and parses a JWT token
func validateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &AuthError{Message: "unexpected signing method"}
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, &AuthError{Message: "invalid token", Err: err}
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, &AuthError{Message: "invalid token claims"}
}

// shouldSkipAuth checks if the path should skip authentication
func shouldSkipAuth(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// handleAuthError handles authentication errors
func handleAuthError(w http.ResponseWriter, r *http.Request, err error, errorFunc func(http.ResponseWriter, *http.Request, error)) {
	if errorFunc != nil {
		errorFunc(w, r, err)
		return
	}

	// Default error handler
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"unauthorized","message":"` + err.Error() + `"}`))
}

// AuthError represents an authentication error
type AuthError struct {
	Message string
	Err     error
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// GetUserID extracts user ID from request context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserEmail extracts user email from request context
func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// MustGetUserID extracts user ID from context or panics
func MustGetUserID(ctx context.Context) string {
	userID, ok := GetUserID(ctx)
	if !ok {
		panic("user ID not found in context")
	}
	return userID
}
