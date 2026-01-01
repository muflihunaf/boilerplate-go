package handler

import (
	"github.com/muflihunaf/boilerplate-go/internal/service"
)

type Handler struct {
	svc     *service.Service
	authSvc *service.AuthService
}

func New(svc *service.Service, authSvc *service.AuthService) *Handler {
	return &Handler{
		svc:     svc,
		authSvc: authSvc,
	}
}
