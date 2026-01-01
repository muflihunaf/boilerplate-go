package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Application
	Env      string
	Port     string
	LogLevel string

	// Server settings
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// JWT settings
	JWTSecret     string
	JWTExpiration time.Duration
	JWTIssuer     string
}

func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	return &Config{
		Env:           getEnv("APP_ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		ReadTimeout:   getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:  getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:   getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		JWTIssuer:     getEnv("JWT_ISSUER", "boilerplate-go"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return fallback
}

