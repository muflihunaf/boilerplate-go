package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
	"github.com/muflihunaf/boilerplate-go/pkg/response"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

// Auth creates a JWT authentication middleware.
func Auth(jwtService *jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "missing authorization header")
				return
			}

			// Extract Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				switch err {
				case jwt.ErrExpiredToken:
					response.Unauthorized(w, "token has expired")
				default:
					response.Unauthorized(w, "invalid token")
				}
				return
			}

			// Add claims to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from context.
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetEmail extracts email from context.
func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}
