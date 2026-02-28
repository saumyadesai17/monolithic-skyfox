package service

import (
	"context"
	"skyfox/_mocks/repomocks"
	"skyfox/bookings/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRevenueService(t *testing.T) {

	t.Run("RevenueByDate when shows exists", func(t *testing.T) {
		showrepo := repomocks.ShowRepository{}
		bookrepo := repomocks.BookingRepository{}
		shows := make([]model.Show, 3)
		expected := 2564.75 //some random number of type float64
		showrepo.On("GetAllShowsOn", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(shows, nil).Once()
		bookrepo.On("BookingAmountByShows", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("[]int")).Return(expected).Once()

		service := NewRevenueService(&bookrepo, &showrepo)

		got, err := service.RevenueOn(context.Background(), "")

		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("RevenueByDate when no shows exist", func(t *testing.T) {
		showrepo := repomocks.ShowRepository{}
		bookrepo := repomocks.BookingRepository{}
		shows := make([]model.Show, 0)
		expected := float64(0)
		showrepo.On("GetAllShowsOn", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(shows, nil).Once()
		bookrepo.On("BookingAmountByShows", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("[]int")).Return(expected).Once()

		service := NewRevenueService(&bookrepo, &showrepo)

		got, err := service.RevenueOn(context.Background(), "")

		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})
}
