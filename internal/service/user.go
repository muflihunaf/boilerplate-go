package service

import (
	"context"

	"github.com/muflihunaf/boilerplate-go/internal/repository"
)

func (s *Service) ListUsers(ctx context.Context) ([]repository.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) GetUser(ctx context.Context, id string) (*repository.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, name, email string) (*repository.User, error) {
	// Add business logic here (e.g., validation, uniqueness check)
	return s.repo.CreateUser(ctx, name, email)
}

func (s *Service) UpdateUser(ctx context.Context, id, name, email string) (*repository.User, error) {
	user, err := s.repo.UpdateUser(ctx, id, name, email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.DeleteUser(ctx, id)
	if err == repository.ErrNotFound {
		return ErrNotFound
	}
	return err
}
