package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/muflihunaf/boilerplate-go/internal/service"
)

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		InternalError(w)
		return
	}

	JSON(w, http.StatusOK, users)
}

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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	// TODO: Add validation with validator package

	user, err := h.svc.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		InternalError(w)
		return
	}

	JSON(w, http.StatusCreated, user)
}

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

