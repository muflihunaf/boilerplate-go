package handler

import (
	"net/http"

	"github.com/muflihunaf/boilerplate-go/internal/middleware"
)

type MeResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// Me returns the current authenticated user's information.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		Unauthorized(w, "user not found in context")
		return
	}

	email, _ := middleware.GetEmail(r.Context())

	JSON(w, http.StatusOK, MeResponse{
		UserID: userID,
		Email:  email,
	})
}

