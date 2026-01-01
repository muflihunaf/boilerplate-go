package jwt_test

import (
	"testing"
	"time"

	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

func TestGenerateAndValidateToken(t *testing.T) {
	svc := jwt.NewService(jwt.Config{
		Secret:     "test-secret-key-32-chars-long!!",
		Expiration: time.Hour,
		Issuer:     "test-issuer",
	})

	// Generate token
	token, err := svc.GenerateToken("user-123", "test@example.com")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	// Validate token
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected user ID 'user-123', got '%s'", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", claims.Email)
	}
}

func TestInvalidToken(t *testing.T) {
	svc := jwt.NewService(jwt.Config{
		Secret:     "test-secret",
		Expiration: time.Hour,
		Issuer:     "test",
	})

	_, err := svc.ValidateToken("invalid-token")
	if err != jwt.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestExpiredToken(t *testing.T) {
	svc := jwt.NewService(jwt.Config{
		Secret:     "test-secret",
		Expiration: -time.Hour, // Already expired
		Issuer:     "test",
	})

	token, err := svc.GenerateToken("user-123", "test@example.com")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err != jwt.ErrExpiredToken {
		t.Errorf("expected ErrExpiredToken, got %v", err)
	}
}

