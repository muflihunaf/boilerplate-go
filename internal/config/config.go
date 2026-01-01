package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	// App
	Env      string
	Port     string
	LogLevel string

	// Server
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration
	JWTIssuer     string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load() // Ignore error - .env is optional

	cfg := &Config{
		Env:           env("APP_ENV", "development"),
		Port:          env("PORT", "8080"),
		LogLevel:      env("LOG_LEVEL", "info"),
		ReadTimeout:   duration("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:  duration("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:   duration("IDLE_TIMEOUT", 60*time.Second),
		JWTSecret:     env("JWT_SECRET", ""),
		JWTExpiration: duration("JWT_EXPIRATION", 24*time.Hour),
		JWTIssuer:     env("JWT_ISSUER", "boilerplate-go"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsProd returns true if running in production.
func (c *Config) IsProd() bool {
	return c.Env == "production"
}

func (c *Config) validate() error {
	// JWT secret is required in production
	if c.IsProd() {
		if c.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET is required in production")
		}
		if len(c.JWTSecret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters")
		}
	}

	// Default secret for development only
	if c.JWTSecret == "" {
		c.JWTSecret = "dev-secret-do-not-use-in-production"
	}

	return nil
}

// --- Helpers ---

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if sec, err := strconv.Atoi(v); err == nil {
			return time.Duration(sec) * time.Second
		}
	}
	return fallback
}
