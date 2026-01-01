package handler

import (
	"encoding/json"
	"net/http"

	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// Login authenticates a user and returns a JWT token.
// This is a simplified example - in production, verify against your user store.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		BadRequest(w, "email and password are required")
		return
	}

	// TODO: In production, verify credentials against your user store
	// Example:
	//   user, err := h.svc.AuthenticateUser(r.Context(), req.Email, req.Password)
	//   if err != nil {
	//       Unauthorized(w, "invalid credentials")
	//       return
	//   }

	// For demo purposes, accept any email/password
	// Replace with actual authentication logic
	userID := "demo-user-id"

	token, err := h.jwt.GenerateToken(userID, req.Email)
	if err != nil {
		InternalError(w)
		return
	}

	JSON(w, http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresIn: 86400, // 24 hours in seconds
	})
}

// AuthHandler holds auth-specific dependencies.
type AuthHandler struct {
	jwt *jwt.Service
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(jwtService *jwt.Service) *AuthHandler {
	return &AuthHandler{
		jwt: jwtService,
	}
}

// Login authenticates and returns a JWT token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		BadRequest(w, "email and password are required")
		return
	}

	// TODO: Verify credentials against your user store
	userID := "demo-user-id"

	token, err := h.jwt.GenerateToken(userID, req.Email)
	if err != nil {
		InternalError(w)
		return
	}

	JSON(w, http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresIn: 86400,
	})
}

