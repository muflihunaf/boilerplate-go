package handler

import (
	"encoding/json"
	"net/http"

	"github.com/muflihunaf/boilerplate-go/internal/middleware"
	"github.com/muflihunaf/boilerplate-go/internal/service"
)

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// AuthResponse represents the authentication response.
type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expires_in"`
	User      UserResponse `json:"user"`
}

// UserResponse represents user data in responses.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Login authenticates a user and returns a JWT token.
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

	result, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			Unauthorized(w, "invalid email or password")
		default:
			InternalError(w)
		}
		return
	}

	JSON(w, http.StatusOK, AuthResponse{
		Token:     result.Token,
		ExpiresIn: 86400, // 24 hours in seconds
		User: UserResponse{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
		},
	})
}

// Register creates a new user and returns a JWT token.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		BadRequest(w, "name, email and password are required")
		return
	}

	if len(req.Password) < 6 {
		BadRequest(w, "password must be at least 6 characters")
		return
	}

	result, err := h.authSvc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrConflict:
			Error(w, http.StatusConflict, "CONFLICT", "email already registered")
		default:
			InternalError(w)
		}
		return
	}

	JSON(w, http.StatusCreated, AuthResponse{
		Token:     result.Token,
		ExpiresIn: 86400,
		User: UserResponse{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
		},
	})
}

// Me returns the current authenticated user's information.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		Unauthorized(w, "user not found in context")
		return
	}

	user, err := h.authSvc.GetCurrentUser(r.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			NotFound(w, "user not found")
		default:
			InternalError(w)
		}
		return
	}

	JSON(w, http.StatusOK, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
