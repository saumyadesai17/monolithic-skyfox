package repository

import (
	"context"
	"errors"
	"fmt"
	"skyfox/bookings/database/common"
	"skyfox/bookings/model"
	ae "skyfox/error"

	"gorm.io/gorm"
)

type UserAccountRepository interface {
	FindByPhone(ctx context.Context, phone string) (*model.UserAccount, error)
	FindByEmail(ctx context.Context, phone string) (*model.UserAccount, error)
	CreateUser(ctx context.Context, user *model.UserAccount) error
}

type userAccountRepository struct {
	*common.BaseDB
}

func NewAccountRepository(db *common.BaseDB) UserAccountRepository {
	return &userAccountRepository{
		BaseDB: db,
	}
}

func (r *userAccountRepository) FindByPhone(ctx context.Context, phone string) (*model.UserAccount, error) {
	var user model.UserAccount

	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // not found is not an error here, caller checks nil
		}
		return nil, ae.InternalServerError("InternalServerError", "something went wrong", fmt.Errorf("FindByPhone failed: %w", result.Error))
	}

	return &user, nil
}

func (r *userAccountRepository) FindByEmail(ctx context.Context, email string) (*model.UserAccount, error) {
	var user model.UserAccount

	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // not found is not an error here, caller checks nil
		}
		return nil, ae.InternalServerError("InternalServerError", "something went wrong", fmt.Errorf("FindByPhone failed: %w", result.Error))
	}

	return &user, nil
}

func (r *userAccountRepository) CreateUser(ctx context.Context, user *model.UserAccount) error {
	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Create(user)
	if result.Error != nil {
		return ae.InternalServerError("InternalServerError", "failed to create user", fmt.Errorf("CreateUser failed: %w", result.Error))
	}

	return nil
}
