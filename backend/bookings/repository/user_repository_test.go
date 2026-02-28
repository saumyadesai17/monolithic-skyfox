package repository_test

import (
	"context"
	"reflect"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/integration_test/db"
	"testing"

	_ "github.com/lib/pq"
)

func TestUserRepository(t *testing.T) {
	//container and database
	db := db.GetDB()
	gormDB := db.GormDB()

	//migrate
	err := gormDB.AutoMigrate(model.User{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	repo := repository.NewUserRepository(db)

	//tests
	t.Run("SaveUser", func(t *testing.T) {
		user := model.NewUser("hello", "hi")
		err = repo.Create(ctx, &user)
		if err != nil {
			t.Errorf("failed to create user: %s", err)
		}
	})

	t.Run("FindExistingUserByUsername", func(t *testing.T) {
		user := model.NewUser("hello", "hi")
		err := repo.Create(ctx, &user)
		if err != nil {
			t.Fatal(err)
		}

		expected := model.User{Id: 1, Username: "hello", Password: "hi"}
		actualUser, err := repo.FindByUsername(ctx, "hello")

		if err != nil {
			t.Fatal(err)
		}

		expected.Id = actualUser.Id // Ignore Id as it auto-increments

		if !(reflect.DeepEqual(expected, actualUser)) {
			t.Errorf("failed to find user: %v", actualUser.Id)
		}
	})
}
