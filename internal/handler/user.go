package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/muflihunaf/boilerplate-go/internal/service"
)

// CreateUserRequest represents the request to create a user.
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// UpdateUserRequest represents the request to update a user.
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100" example:"John Doe"`
	Email string `json:"email" validate:"omitempty,email" example:"user@example.com"`
}

// ListUsers godoc
// @Summary      List all users
// @Description  Returns a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   UserResponse
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		InternalError(w)
		return
	}

	JSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Returns a single user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  UserResponse
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		BadRequest(w, "id is required")
		return
	}

	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			NotFound(w, "user not found")
			return
		}
		InternalError(w)
		return
	}

	JSON(w, http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateUserRequest  true  "User details"
// @Success      201      {object}  UserResponse
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	user, err := h.svc.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		InternalError(w)
		return
	}

	JSON(w, http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Updates an existing user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string             true  "User ID"
// @Param        request  body      UpdateUserRequest  true  "User details"
// @Success      200      {object}  UserResponse
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      404      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		BadRequest(w, "id is required")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	user, err := h.svc.UpdateUser(r.Context(), id, req.Name, req.Email)
	if err != nil {
		if err == service.ErrNotFound {
			NotFound(w, "user not found")
			return
		}
		InternalError(w)
		return
	}

	JSON(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      204  "No Content"
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		BadRequest(w, "id is required")
		return
	}

	if err := h.svc.DeleteUser(r.Context(), id); err != nil {
		if err == service.ErrNotFound {
			NotFound(w, "user not found")
			return
		}
		InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
