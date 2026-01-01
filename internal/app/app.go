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
)

// App holds all application dependencies.
type App struct {
	cfg    *config.Config
	logger *slog.Logger
	server *server.Server
}

// New creates a new application instance.
func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	logger := initLogger(cfg)

	// Initialize layers (dependency injection)
	repo := repository.New()
	svc := service.New(repo)
	h := handler.New(svc)

	// Create server
	srv := server.New(cfg, h, logger)

	return &App{
		cfg:    cfg,
		logger: logger,
		server: srv,
	}, nil
}

// Run starts the application and blocks until shutdown.
func (a *App) Run() error {
	// Start server
	errChan := make(chan error, 1)
	go func() {
		a.logger.Info("starting server",
			"port", a.cfg.Port,
			"env", a.cfg.Env,
		)
		if err := a.server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	return a.waitForShutdown(errChan)
}

func (a *App) waitForShutdown(errChan chan error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		a.logger.Info("received shutdown signal", "signal", sig.String())
	}

	return a.shutdown()
}

func (a *App) shutdown() error {
	a.logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server forced to shutdown", "error", err)
		return err
	}

	a.logger.Info("server exited gracefully")
	return nil
}

func initLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}

	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

