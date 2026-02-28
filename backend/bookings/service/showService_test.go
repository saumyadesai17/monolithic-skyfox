package service

import (
	"context"
	"fmt"
	"skyfox/_mocks/repomocks"
	"skyfox/bookings/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShowService(t *testing.T) {

	t.Run("GetShows", func(t *testing.T) {
		showrepo := repomocks.ShowRepository{}
		movieGatewayRepo := repomocks.MovieGateWay{}
		expected := make([]model.Show, 0)
		showrepo.On("GetAllShowsOn", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(expected, nil).Once()

		service := NewShowService(&showrepo, &movieGatewayRepo)

		got, err := service.GetShows(context.Background(), "")

		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("GetMovieById when error is nil", func(t *testing.T) {
		showRepo := repomocks.ShowRepository{}
		movieGatewayRepo := repomocks.MovieGateWay{}
		expected := &model.Movie{
			MovieId:  "id",
			Name:     "movie",
			Duration: "1h30m",
			Plot:     "This is a horror movie.",
		}
		movieGatewayRepo.On("MovieById", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(expected, nil).Once()

		service := NewShowService(&showRepo, &movieGatewayRepo)

		got, err := service.GetMovieById(context.Background(), "")

		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("GetMovieById when error is not nil", func(t *testing.T) {
		showRepo := repomocks.ShowRepository{}
		movieGatewayRepo := repomocks.MovieGateWay{}
		expected := &model.Movie{}
		someError := fmt.Errorf("some error")
		movieGatewayRepo.On("MovieById", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(expected, someError).Once()

		service := NewShowService(&showRepo, &movieGatewayRepo)

		got, err := service.GetMovieById(context.Background(), "")

		assert.NotNil(t, err)
		assert.Equal(t, expected, got)
	})

}
