package account

import (
	"context"
	"strings"

	"mini-payment-system/internal/domain/entities"
	"mini-payment-system/internal/domain/errs"
	"mini-payment-system/internal/domain/repositories"

	"github.com/google/uuid"
)

type Service struct {
	repo     repositories.AccountRepository
	userRepo repositories.UserRepository
}

func NewService(repo repositories.AccountRepository, userRepo repositories.UserRepository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) Create(ctx context.Context, userID string, initialBalance int64, currency string) (*entities.Account, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errs.InvalidInput("user_id is required")
	}
	if initialBalance < 0 {
		return nil, errs.InvalidInput("initial balance cannot be negative")
	}
	if strings.TrimSpace(currency) == "" {
		currency = "THB"
	}

	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, errs.NotFound("user not found")
	}

	account := &entities.Account{
		ID:       uuid.NewString(),
		UserID:   userID,
		Balance:  initialBalance,
		Currency: strings.ToUpper(strings.TrimSpace(currency)),
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*entities.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]entities.Account, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id, currency string) (*entities.Account, error) {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(currency) != "" {
		account.Currency = strings.ToUpper(strings.TrimSpace(currency))
	}

	if err := s.repo.Update(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
