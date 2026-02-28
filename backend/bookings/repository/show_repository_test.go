package repository_test

import (
	"context"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/bookings/repository/testdata"
	"skyfox/integration_test/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowRepository(t *testing.T) {
	//container and database
	db := db.GetDB()

	//migrate
	err := db.GormDB().AutoMigrate(model.Show{}, model.Slot{})
	db.GormDB().Exec("DELETE FROM BOOKING")
	db.GormDB().Exec("DELETE FROM SHOW")
	db.GormDB().Exec("DELETE FROM SLOT")
	db.GormDB().Create(testdata.Shows)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	repo := repository.NewShowRepository(db)

	//tests
	t.Run("FindShowById", func(t *testing.T) {
		expected := model.Show{
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
		}

		actual, _ := repo.FindById(ctx, 1)

		assert.Equal(t, expected, actual)
	})

	t.Run("GetAllShowsByDate", func(t *testing.T) {
		expected := []model.Show{
			{
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
			},
			{
				Id:      2,
				MovieId: "tt6856489",
				Date:    "2022-10-13",
				SlotId:  4,
				Slot: model.Slot{
					Id:        4,
					Name:      "slot4",
					StartTime: "22:30:00",
					EndTime:   "02:00:00",
				},
				Cost: 350.00,
			},
			{
				Id:      3,
				MovieId: "tt6856999",
				Date:    "2022-10-13",
				SlotId:  1,
				Slot: model.Slot{
					Id:        1,
					Name:      "slot1",
					StartTime: "09:00:00",
					EndTime:   "12:30:00",
				},
				Cost: 350.00,
			},
		}

		actual, err := repo.GetAllShowsOn(ctx, "2022-10-13")

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
