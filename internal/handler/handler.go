package handler

import (
	"github.com/muflihunaf/boilerplate-go/internal/service"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

type Handler struct {
	svc *service.Service
	jwt *jwt.Service
}

func New(svc *service.Service, jwtService *jwt.Service) *Handler {
	return &Handler{
		svc: svc,
		jwt: jwtService,
	}
}

