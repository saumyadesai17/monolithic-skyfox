package repository

import (
	// "skyfox/bookings/database/connection"
	"context"
	"errors"
	"fmt"
	"skyfox/bookings/database/common"
	"skyfox/bookings/model"
	ae "skyfox/error"

	"gorm.io/gorm"
)

type userRepository struct {
	*common.BaseDB
}

func NewUserRepository(db *common.BaseDB) *userRepository {
	return &userRepository{
		BaseDB: db,
	}
}

func (repo userRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	db, cancel := repo.WithContext(ctx)
	defer cancel()

	if result := db.Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.User{}, nil
		}
		return model.User{}, ae.InternalServerError("InternalServerError", fmt.Sprintf("something went wrong"),
			fmt.Errorf("something went wrong %v", result.Error))
	}

	return user, nil
}

func (repo userRepository) Create(ctx context.Context, user *model.User) error {

	db, cancel := repo.WithContext(ctx)
	defer cancel()

	result := db.Create(user)

	if result.Error != nil {
		return ae.UnProcessableError("UserCreationFailed", "User creation failed due to unknown reason", result.Error)
	}
	return nil
}
