package service

import (
	"context"
	"skyfox/bookings/model"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepository UserRepository) *userService {
	return &userService{
		userRepo: userRepository,
	}
}

func (s *userService) UserDetails(ctx context.Context, username string) (model.User, error) {

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
