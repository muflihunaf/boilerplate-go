package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muflihunaf/boilerplate-go/internal/handler"
	"github.com/muflihunaf/boilerplate-go/internal/repository"
	"github.com/muflihunaf/boilerplate-go/internal/service"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

func TestHealth(t *testing.T) {
	// Setup
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Execute
	h.Health(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Status  string `json:"status"`
			Version string `json:"version"`
		} `json:"data"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.Data.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Data.Status)
	}
}

func newTestHandler() *handler.Handler {
	repo := repository.New()
	svc := service.New(repo)
	jwtSvc := jwt.NewService(jwt.Config{
		Secret:     "test-secret",
		Expiration: 3600,
		Issuer:     "test",
	})
	authSvc := service.NewAuthService(repo, jwtSvc, 3600)
	return handler.New(svc, authSvc)
}

