package repositories

import (
	"context"

	"mini-payment-system/internal/domain/entities"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entities.Account) error
	GetByID(ctx context.Context, id string) (*entities.Account, error)
	GetByIDForUpdate(ctx context.Context, id string) (*entities.Account, error)
	List(ctx context.Context) ([]entities.Account, error)
	Update(ctx context.Context, account *entities.Account) error
	Delete(ctx context.Context, id string) error
}
