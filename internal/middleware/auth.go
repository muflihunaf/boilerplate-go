package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/muflihunaf/boilerplate-go/internal/handler"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
)

// Auth is a middleware that validates authentication tokens.
// Replace with your actual auth logic (JWT, session, API key, etc.)
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handler.Unauthorized(w, "missing authorization header")
			return
		}

		// Example: Bearer token parsing
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.Unauthorized(w, "invalid authorization header format")
			return
		}

		token := parts[1]

		// TODO: Validate token and extract user ID
		// This is a placeholder - implement your actual token validation
		userID, err := validateToken(token)
		if err != nil {
			handler.Unauthorized(w, "invalid token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// validateToken is a placeholder for actual token validation
func validateToken(token string) (string, error) {
	// TODO: Implement actual JWT/token validation
	// Example with JWT:
	//   claims, err := jwt.Parse(token, secret)
	//   return claims.UserID, err

	// For now, just return the token as user ID (placeholder)
	if token == "" {
		return "", http.ErrNoCookie
	}
	return token, nil
}

