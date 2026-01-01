package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muflihunaf/boilerplate-go/internal/config"
	"github.com/muflihunaf/boilerplate-go/internal/handler"
	"github.com/muflihunaf/boilerplate-go/internal/repository"
	"github.com/muflihunaf/boilerplate-go/internal/server"
	"github.com/muflihunaf/boilerplate-go/internal/service"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

// App holds all application dependencies.
type App struct {
	cfg    *config.Config
	log    *slog.Logger
	server *server.Server
}

// New creates a new application instance.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := setupLogger(cfg)
	jwtSvc := jwt.NewService(jwt.Config{
		Secret:     cfg.JWTSecret,
		Expiration: cfg.JWTExpiration,
		Issuer:     cfg.JWTIssuer,
	})

	// Wire dependencies
	repo := repository.New()
	svc := service.New(repo)
	authSvc := service.NewAuthService(repo, jwtSvc, cfg.JWTExpiration)
	h := handler.New(svc, authSvc)

	return &App{
		cfg:    cfg,
		log:    log,
		server: server.New(cfg, h, jwtSvc, log),
	}, nil
}

// Run starts the server and blocks until shutdown.
func (a *App) Run() error {
	errCh := make(chan error, 1)
	go func() {
		a.log.Info("starting server", "port", a.cfg.Port, "env", a.cfg.Env)
		if err := a.server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	return a.awaitShutdown(errCh)
}

func (a *App) awaitShutdown(errCh chan error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-quit:
		a.log.Info("received signal", "signal", sig.String())
	}

	a.log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("shutdown error", "error", err)
		return err
	}

	a.log.Info("server stopped")
	return nil
}

func setupLogger(cfg *config.Config) *slog.Logger {
	level := parseLevel(cfg.LogLevel)
	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if cfg.IsProd() {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
