package repository_test

import (
	"context"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/bookings/repository/testdata"
	"skyfox/integration_test/db"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestBookingRepository(t *testing.T) {
	//container and database
	db := db.GetDB()

	//migrate
	err := db.GormDB().AutoMigrate(model.Booking{}, model.Show{}, model.Customer{}, model.Slot{})

	db.GormDB().Exec("DELETE FROM BOOKING")
	db.GormDB().Exec("DELETE FROM SHOW")
	db.GormDB().Exec("DELETE FROM SLOT")
	db.GormDB().Create(testdata.Shows)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	repo := repository.NewBookingRepository(db)

	//tests
	t.Run("SaveBooking", func(t *testing.T) {
		booking := model.NewBooking("2022-10-13", model.Show{
			Id:      1,
			MovieId: "tt6857189",
			Date:    "2022-10-13",
			SlotId:  3,
			Slot: model.Slot{
				Id:        3,
				Name:      "slot3",
				StartTime: "18:00:00",
				EndTime:   "21:30:00",
			},
			Cost: 300.00,
		}, model.Customer{Id: 1, Name: "John", PhoneNumber: "6543276543"}, 2, 600.00)

		err = repo.Create(ctx, &booking)
		if err != nil {
			t.Errorf("failed to book: %s", err)
		}
	})

	t.Run("FindBookedSeatsByShow", func(t *testing.T) {
		booking := model.NewBooking("2022-10-13", model.Show{
			Id:      1,
			MovieId: "tt6857189",
			Date:    "2022-10-13",
			SlotId:  3,
			Slot: model.Slot{
				Id:        3,
				Name:      "slot3",
				StartTime: "18:00:00",
				EndTime:   "21:30:00",
			},
			Cost: 300.00,
		}, model.Customer{Id: 1, Name: "John", PhoneNumber: "6543276543"}, 2, 600.00)
		repo.Create(ctx, &booking)
		repo.Create(ctx, &booking)

		expectedSeats := int(4)

		actualSeats := repo.BookedSeatsByShow(ctx, 1)

		assert.Equal(t, expectedSeats, actualSeats)
	})
}
