package persistence

import (
	"context"
	"skyfox/bookings/model"
	"skyfox/bookings/service"
)

func SeedDB(userRepo service.UserRepository) {

	ctx := context.Background()
	user, _ := userRepo.FindByUsername(ctx, "seed-user-1")
	if user == (model.User{}) {
		user := model.NewUser("seed-user-1", "foobar")
		userRepo.Create(ctx, &user)

	}
	user, _ = userRepo.FindByUsername(ctx, "seed-user-2")
	if user == (model.User{}) {
		user := model.NewUser("seed-user-2", "foobar")
		userRepo.Create(ctx, &user)
	}
}
