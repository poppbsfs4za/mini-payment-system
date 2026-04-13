package gormrepo

import (
	"context"

	"mini-payment-system/internal/domain/entities"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	return getDB(ctx, r.db).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	if err := getDB(ctx, r.db).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	if err := getDB(ctx, r.db).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	return getDB(ctx, r.db).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return getDB(ctx, r.db).Delete(&entities.User{}, "id = ?", id).Error
}
