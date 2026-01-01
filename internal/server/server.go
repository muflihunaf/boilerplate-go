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

// Server represents the HTTP server.
type Server struct {
	httpServer *http.Server
	router     *chi.Mux
	logger     *slog.Logger
}

// New creates a new HTTP server with all routes and middleware configured.
func New(cfg *config.Config, h *handler.Handler, jwtService *jwt.Service, logger *slog.Logger) *Server {
	r := chi.NewRouter()

	// Setup global middleware stack
	SetupMiddleware(r)

	// Register all routes
	RegisterRoutes(r, h, jwtService)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		router: r,
		logger: logger,
	}
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Router returns the underlying chi router (for testing).
func (s *Server) Router() *chi.Mux {
	return s.router
}
