package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/muflihunaf/boilerplate-go/internal/config"
	"github.com/muflihunaf/boilerplate-go/internal/handler"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
	logger     *slog.Logger
}

func New(cfg *config.Config, h *handler.Handler, logger *slog.Logger) *Server {
	r := chi.NewRouter()

	// Production-ready middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware (configure as needed)
	r.Use(corsMiddleware)

	// Register routes
	registerRoutes(r, h)

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

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func registerRoutes(r *chi.Mux, h *handler.Handler) {
	// Health check (public)
	r.Get("/health", h.Health)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Example: Users resource
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUsers)
			r.Post("/", h.CreateUser)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetUser)
				r.Put("/", h.UpdateUser)
				r.Delete("/", h.DeleteUser)
			})
		})

		// Add more resources here
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

