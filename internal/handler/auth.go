package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/muflihunaf/boilerplate-go/internal/middleware"
	"github.com/muflihunaf/boilerplate-go/internal/service"
)

// --- Request/Response Types ---

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expires_in"`
	User      UserResponse `json:"user"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// --- Handlers ---

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200      {object}  AuthResponse
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Router       /auth/login [post]
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
		if err == service.ErrInvalidCredentials {
			Unauthorized(w, "invalid email or password")
			return
		}
		slog.Error("login failed", "error", err, "email", req.Email)
		InternalError(w)
		return
	}

	OK(w, toAuthResponse(result))
}

// Register godoc
// @Summary      User registration
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration details"
// @Success      201      {object}  AuthResponse
// @Failure      400      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Router       /auth/register [post]
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
		if err == service.ErrConflict {
			Conflict(w, "email already registered")
			return
		}
		slog.Error("registration failed", "error", err, "email", req.Email)
		InternalError(w)
		return
	}

	Created(w, toAuthResponse(result))
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's profile
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UserResponse
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		Unauthorized(w, "user not found in context")
		return
	}

	user, err := h.authSvc.GetCurrentUser(r.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			NotFound(w, "user not found")
			return
		}
		slog.Error("get current user failed", "error", err, "user_id", userID)
		InternalError(w)
		return
	}

	OK(w, UserResponse{ID: user.ID, Name: user.Name, Email: user.Email})
}

// --- Helpers ---

func toAuthResponse(r *service.AuthResult) AuthResponse {
	return AuthResponse{
		Token:     r.Token,
		ExpiresIn: 86400,
		User:      UserResponse{ID: r.User.ID, Name: r.User.Name, Email: r.User.Email},
	}
}

