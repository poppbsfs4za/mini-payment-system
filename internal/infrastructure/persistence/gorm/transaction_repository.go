package gormrepo

import (
	"context"

	"mini-payment-system/internal/domain/entities"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *entities.Transaction) error {
	return getDB(ctx, r.db).Create(transaction).Error
}

func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := getDB(ctx, r.db).First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) List(ctx context.Context) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	if err := getDB(ctx, r.db).Order("created_at desc").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
