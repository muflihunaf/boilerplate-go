package service

import (
	"context"
	"errors"
	"time"

	"github.com/muflihunaf/boilerplate-go/internal/repository"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthResult contains the authentication result.
type AuthResult struct {
	Token     string
	ExpiresAt time.Time
	User      *repository.User
}

// AuthService handles authentication logic.
type AuthService struct {
	repo       *repository.Repository
	jwtService *jwt.Service
	expiration time.Duration
}

// NewAuthService creates a new auth service.
func NewAuthService(repo *repository.Repository, jwtService *jwt.Service, expiration time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
		expiration: expiration,
	}
}

// Login authenticates a user and returns a token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	// Find user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if !s.repo.CheckPassword(user.Password, password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token:     token,
		ExpiresAt: time.Now().Add(s.expiration),
		User:      user,
	}, nil
}

// Register creates a new user and returns a token.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*AuthResult, error) {
	// Check if user already exists
	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, ErrConflict
	}

	// Create user with hashed password
	user, err := s.repo.CreateUserWithPassword(ctx, name, email, password)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token:     token,
		ExpiresAt: time.Now().Add(s.expiration),
		User:      user,
	}, nil
}

// GetCurrentUser returns the user for a given user ID.
func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*repository.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

