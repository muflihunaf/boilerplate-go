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

// Auth validates JWT tokens and injects user claims into context.
func Auth(jwtSvc *jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				response.Unauthorized(w, "missing authorization header")
				return
			}

			claims, err := jwtSvc.ValidateToken(token)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					response.Unauthorized(w, "token has expired")
				} else {
					response.Unauthorized(w, "invalid token")
				}
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// GetUserID extracts user ID from context.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

// GetEmail extracts email from context.
func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}
