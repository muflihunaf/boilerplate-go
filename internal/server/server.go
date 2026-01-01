package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/muflihunaf/boilerplate-go/internal/config"
	"github.com/muflihunaf/boilerplate-go/internal/handler"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

// Server wraps the HTTP server.
type Server struct {
	http   *http.Server
	router *chi.Mux
	log    *slog.Logger
}

// New creates a configured HTTP server.
func New(cfg *config.Config, h *handler.Handler, jwtSvc *jwt.Service, log *slog.Logger) *Server {
	r := chi.NewRouter()
	SetupMiddleware(r)
	RegisterRoutes(r, h, jwtSvc)

	return &Server{
		http: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		router: r,
		log:    log,
	}
}

// Start begins listening for requests.
func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

// Router returns the chi router for testing.
func (s *Server) Router() *chi.Mux {
	return s.router
}
