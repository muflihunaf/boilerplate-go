package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/muflihunaf/boilerplate-go/internal/handler"
	"github.com/muflihunaf/boilerplate-go/internal/middleware"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

// RegisterRoutes sets up all application routes.
// No business logic here - only route definitions.
func RegisterRoutes(r *chi.Mux, h *handler.Handler, jwtService *jwt.Service) {
	// Public routes (no authentication required)
	registerPublicRoutes(r, h)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public API routes
		registerAuthRoutes(r, h)

		// Protected API routes (require JWT)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtService))

			registerUserRoutes(r, h)
			registerProfileRoutes(r, h)
		})
	})
}

// registerPublicRoutes sets up routes that don't require authentication.
func registerPublicRoutes(r chi.Router, h *handler.Handler) {
	// Health check
	r.Get("/health", h.Health)

	// Ready check (can add DB ping, etc.)
	r.Get("/ready", h.Health)
}

// registerAuthRoutes sets up authentication routes.
func registerAuthRoutes(r chi.Router, h *handler.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})
}

// registerUserRoutes sets up user management routes (protected).
func registerUserRoutes(r chi.Router, h *handler.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Post("/", h.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUser)
			r.Put("/", h.UpdateUser)
			r.Delete("/", h.DeleteUser)
		})
	})
}

// registerProfileRoutes sets up current user profile routes (protected).
func registerProfileRoutes(r chi.Router, h *handler.Handler) {
	r.Get("/me", h.Me)
}

