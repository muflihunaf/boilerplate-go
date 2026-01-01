package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// SetupMiddleware configures the global middleware stack.
// Order matters: middleware is executed in the order it's added.
func SetupMiddleware(r interface{ Use(middlewares ...func(http.Handler) http.Handler) }) {
	// Request ID - assigns unique ID to each request
	r.Use(middleware.RequestID)

	// Real IP - extracts real client IP from proxy headers
	r.Use(middleware.RealIP)

	// Logger - logs request details
	r.Use(middleware.Logger)

	// Recoverer - recovers from panics and returns 500
	r.Use(middleware.Recoverer)

	// CleanPath - cleans double slashes from URL paths
	r.Use(middleware.CleanPath)

	// Timeout - cancels context after specified duration
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS - handles Cross-Origin Resource Sharing
	r.Use(CORS)

	// Security headers
	r.Use(SecureHeaders)
}

// CORS middleware handles Cross-Origin Resource Sharing.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecureHeaders adds security-related headers to responses.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS filter
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// ContentType sets the Content-Type header for JSON responses.
func ContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}

