package repository_test

import (
	"context"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/integration_test/db"
	"testing"

	_ "github.com/lib/pq"
)

func TestCustomerRepository(t *testing.T) {
	//container and database
	db := db.GetDB()

	//migrate
	err := db.GormDB().AutoMigrate(model.Customer{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	repo := repository.NewCustomerRepository(db)

	//tests
	t.Run("SaveCustomer", func(t *testing.T) {
		customer := model.Customer{Id: 1, Name: "John", PhoneNumber: "6543276543"}
		err = repo.Create(ctx, &customer)
		if err != nil {
			t.Errorf("failed to create customer: %s", err)
		}
	})
}
