package repositories

import (
	"context"

	"mini-payment-system/internal/domain/entities"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	GetByID(ctx context.Context, id string) (*entities.Transaction, error)
	List(ctx context.Context) ([]entities.Transaction, error)
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
