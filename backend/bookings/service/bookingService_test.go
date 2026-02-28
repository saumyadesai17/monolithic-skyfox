package service

import (
	"context"
	"skyfox/_mocks/repomocks"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/model"
	ae "skyfox/error"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingService(t *testing.T) {

	var bookedSeats int
	show := model.Show{
		Id:      0,
		MovieId: "movie_id",
		Date:    "2022-11-19",
		Slot: model.Slot{
			Id:        3,
			Name:      "slot3",
			StartTime: "5:00:00",
			EndTime:   "8:00:00",
		},
		SlotId: 1,
		Cost:   350.64,
	}
	want := &model.Booking{
		Id:     0,
		Date:   "",
		Show:   show,
		ShowId: 0,
		Customer: model.Customer{
			Id:          0,
			Name:        "abc",
			PhoneNumber: "9876543210",
		},
		CustomerId: 0,
		NoOfSeats:  0,
		AmountPaid: 0,
	}
	bookingRequest := request.BookingRequest{
		Date:   "",
		ShowId: 0,
		Customer: model.Customer{
			Name:        "abc",
			PhoneNumber: "9876543210",
		},
		NoOfSeats: 0,
	}

	t.Run("Should book when seats are available", func(t *testing.T) {
		bookrepo := repomocks.BookingRepository{}
		showrepo := repomocks.ShowRepository{}
		bookedSeats = 5
		showrepo.On("FindById", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("int")).Return(show, nil).Once()
		bookrepo.On("BookedSeatsByShow", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("int")).Return(bookedSeats).Once()
		bookrepo.On("Create", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*model.Booking")).Return(nil).Once()

		service := NewBookingService(&bookrepo, &showrepo)

		got, err := service.Book(context.Background(), bookingRequest)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Should not book when seats are not available", func(t *testing.T) {
		bookrepo := repomocks.BookingRepository{}
		showrepo := repomocks.ShowRepository{}
		bookedSeats = 150
		showrepo.On("FindById", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("int")).Return(show, nil).Once()
		bookrepo.On("BookedSeatsByShow", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("int")).Return(bookedSeats).Once()

		service := NewBookingService(&bookrepo, &showrepo)

		got, err := service.Book(context.Background(), bookingRequest)

		assert.Nil(t, got)
		assert.NotNil(t, err)
	})

	t.Run("Should not book when show is not found", func(t *testing.T) {
		bookrepo := repomocks.BookingRepository{}
		showrepo := repomocks.ShowRepository{}
		showrepo.On("FindById", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("int")).Return(model.Show{},
			ae.NotFoundError("Notfound", "mocked", nil)).Once()
		service := NewBookingService(&bookrepo, &showrepo)

		got, showError := service.Book(context.Background(), bookingRequest)
		err := showError.(*ae.AppError)

		assert.Nil(t, got)
		assert.NotNil(t, err)

	})
}
