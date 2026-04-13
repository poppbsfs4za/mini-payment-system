package user

import (
	"context"
	"strings"

	"mini-payment-system/internal/domain/entities"
	"mini-payment-system/internal/domain/errs"
	"mini-payment-system/internal/domain/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo repositories.UserRepository
}

func NewService(repo repositories.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name, email string) (*entities.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if name == "" || email == "" {
		return nil, errs.InvalidInput("name and email are required")
	}

	user := &entities.User{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*entities.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]entities.User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id, name, email string) (*entities.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(name) != "" {
		user.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		user.Email = strings.TrimSpace(strings.ToLower(email))
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
