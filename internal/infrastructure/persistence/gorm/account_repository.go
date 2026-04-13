package gormrepo

import (
	"context"

	"mini-payment-system/internal/domain/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *entities.Account) error {
	return getDB(ctx, r.db).Create(account).Error
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*entities.Account, error) {
	var account entities.Account
	if err := getDB(ctx, r.db).First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, id string) (*entities.Account, error) {
	var account entities.Account
	if err := getDB(ctx, r.db).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) List(ctx context.Context) ([]entities.Account, error) {
	var accounts []entities.Account
	if err := getDB(ctx, r.db).Order("created_at desc").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) Update(ctx context.Context, account *entities.Account) error {
	return getDB(ctx, r.db).Save(account).Error
}

func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	return getDB(ctx, r.db).Delete(&entities.Account{}, "id = ?", id).Error
}
