package server

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/muflihunaf/boilerplate-go/internal/handler"
	"github.com/muflihunaf/boilerplate-go/internal/middleware"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

// RegisterRoutes sets up all application routes.
func RegisterRoutes(r *chi.Mux, h *handler.Handler, jwtSvc *jwt.Service) {
	// Health & docs (public)
	r.Get("/health", h.Health)
	r.Get("/ready", h.Health)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/register", h.Register)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSvc))
			r.Get("/me", h.Me)

			// Users CRUD
			r.Route("/users", func(r chi.Router) {
				r.Get("/", h.ListUsers)
				r.Post("/", h.CreateUser)
				r.Get("/{id}", h.GetUser)
				r.Put("/{id}", h.UpdateUser)
				r.Delete("/{id}", h.DeleteUser)
			})
		})
	})
}
