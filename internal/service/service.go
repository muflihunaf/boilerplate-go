package service

import (
	"errors"

	"github.com/muflihunaf/boilerplate-go/internal/repository"
)

// Common service errors.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = errors.New("resource conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

// Service handles business logic.
type Service struct {
	repo *repository.Repository
}

// New creates a new service.
func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
