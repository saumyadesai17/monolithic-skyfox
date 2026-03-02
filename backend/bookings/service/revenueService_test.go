package service

import (
	"context"
	"skyfox/bookings/model"
	servicemocks "skyfox/bookings/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRevenueService(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository)
		wantRevenue float64
		wantErr     bool
	}{
		{
			name: "Should return total revenue when shows exist on date",
			setupMock: func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository) {
				shows := make([]model.Show, 3)
				showRepo.On("GetAllShowsOn", mock.Anything, mock.AnythingOfType("string")).
					Return(shows, nil).Once()
				bookRepo.On("BookingAmountByShows", mock.Anything, mock.AnythingOfType("[]int")).
					Return(float64(2564.75)).Once()
			},
			wantRevenue: 2564.75,
		},
		{
			name: "Should return zero revenue when no shows exist on date",
			setupMock: func(bookRepo *servicemocks.MockBookingRepository, showRepo *servicemocks.MockShowRepository) {
				showRepo.On("GetAllShowsOn", mock.Anything, mock.AnythingOfType("string")).
					Return([]model.Show{}, nil).Once()
			},
			wantRevenue: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bookRepo := servicemocks.NewMockBookingRepository(t)
			showRepo := servicemocks.NewMockShowRepository(t)
			tc.setupMock(bookRepo, showRepo)

			svc := NewRevenueService(bookRepo, showRepo)
			got, err := svc.RevenueOn(context.Background(), "")

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantRevenue, got)
			}
		})
	}
}
