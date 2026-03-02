package service

import (
	"context"
	"net/http"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/model"
	servicemocks "skyfox/bookings/service/mocks"
	ae "skyfox/error"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingService(t *testing.T) {
	show := model.Show{
		Id:      0,
		MovieId: "movie_id",
		Date:    "2022-11-19",
		Slot:    model.Slot{Id: 3, Name: "slot3", StartTime: "5:00:00", EndTime: "8:00:00"},
		SlotId:  1,
		Cost:    350.64,
	}
	bookingRequest := request.BookingRequest{
		Date:      "",
		ShowId:    0,
		Customer:  model.Customer{Name: "abc", PhoneNumber: "9876543210"},
		NoOfSeats: 0,
	}
	wantBooking := &model.Booking{
		Show:     show,
		Customer: model.Customer{Id: 0, Name: "abc", PhoneNumber: "9876543210"},
	}

	tests := []struct {
		name         string
		setupMock    func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository)
		wantBooking  *model.Booking
		wantHTTPCode int
	}{
		{
			name: "Should book successfully when seats are available",
			setupMock: func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository) {
				showRepo.On("FindById", mock.Anything, mock.AnythingOfType("int")).Return(show, nil).Once()
				bookRepo.On("BookedSeatsByShow", mock.Anything, mock.AnythingOfType("int")).Return(5).Once()
				bookRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Booking")).Return(nil).Once()
			},
			wantBooking: wantBooking,
		},
		{
			name: "Should return 400 when seats are fully booked",
			setupMock: func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository) {
				showRepo.On("FindById", mock.Anything, mock.AnythingOfType("int")).Return(show, nil).Once()
				bookRepo.On("BookedSeatsByShow", mock.Anything, mock.AnythingOfType("int")).Return(150).Once()
			},
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name: "Should return 404 when show does not exist",
			setupMock: func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository) {
				showRepo.On("FindById", mock.Anything, mock.AnythingOfType("int")).
					Return(model.Show{}, ae.NotFoundError("ShowNotFound", "show not found", nil)).Once()
			},
			wantHTTPCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bookRepo := servicemocks.NewMockBookingRepository(t)
			showRepo := servicemocks.NewMockShowRepository(t)
			tc.setupMock(bookRepo, showRepo)

			svc := NewBookingService(bookRepo, showRepo)
			got, err := svc.Book(context.Background(), bookingRequest)

			if tc.wantBooking != nil {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantBooking, got)
			} else {
				assert.Nil(t, got)
				assert.NotNil(t, err)
				appErr := err.(*ae.AppError)
				assert.Equal(t, tc.wantHTTPCode, appErr.HTTPCode())
			}
		})
	}
}
