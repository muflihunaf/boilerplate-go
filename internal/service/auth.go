package service

import (
	"context"
	"time"

	"github.com/muflihunaf/boilerplate-go/internal/repository"
	"github.com/muflihunaf/boilerplate-go/pkg/jwt"
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
	jwt        *jwt.Service
	expiration time.Duration
}

// NewAuthService creates a new auth service.
func NewAuthService(repo *repository.Repository, jwt *jwt.Service, exp time.Duration) *AuthService {
	return &AuthService{repo: repo, jwt: jwt, expiration: exp}
}

// Login authenticates a user and returns a token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !s.repo.CheckPassword(user.Password, password) {
		return nil, ErrInvalidCredentials
	}

	return s.createAuthResult(user)
}

// Register creates a new user and returns a token.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*AuthResult, error) {
	if existing, _ := s.repo.GetUserByEmail(ctx, email); existing != nil {
		return nil, ErrConflict
	}

	user, err := s.repo.CreateUserWithPassword(ctx, name, email, password)
	if err != nil {
		return nil, err
	}

	return s.createAuthResult(user)
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

func (s *AuthService) createAuthResult(user *repository.User) (*AuthResult, error) {
	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token:     token,
		ExpiresAt: time.Now().Add(s.expiration),
		User:      user,
	}, nil
}
