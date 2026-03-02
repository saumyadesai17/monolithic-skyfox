package service

import (
	"context"
	"fmt"
	"skyfox/bookings/model"
	servicemocks "skyfox/bookings/service/mocks"
	moviemocks "skyfox/movieservice/movie_gateway/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShowService_GetShows(t *testing.T) {
	emptyShows := make([]model.Show, 0)
	someShows := []model.Show{{Id: 1, MovieId: "tt0111161"}}

	tests := []struct {
		name      string
		setupMock func(showRepo *servicemocks.MockShowRepository)
		wantShows []model.Show
		wantErr   bool
	}{
		{
			name: "Should return empty list when no shows exist on date",
			setupMock: func(showRepo *servicemocks.MockShowRepository) {
				showRepo.On("GetAllShowsOn", mock.Anything, mock.AnythingOfType("string")).
					Return(emptyShows, nil).Once()
			},
			wantShows: emptyShows,
		},
		{
			name: "Should return shows when they exist for the date",
			setupMock: func(showRepo *servicemocks.MockShowRepository) {
				showRepo.On("GetAllShowsOn", mock.Anything, mock.AnythingOfType("string")).
					Return(someShows, nil).Once()
			},
			wantShows: someShows,
		},
		{
			name: "Should propagate error when repository returns a database error",
			setupMock: func(showRepo *servicemocks.MockShowRepository) {
				showRepo.On("GetAllShowsOn", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, fmt.Errorf("db error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			showRepo := servicemocks.NewMockShowRepository(t)
			gateway := moviemocks.NewMockMovieGateWay(t)
			tc.setupMock(showRepo)

			svc := NewShowService(showRepo, gateway)
			got, err := svc.GetShows(context.Background(), "")

			if tc.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantShows, got)
			}
		})
	}
}

func TestShowService_GetMovieById(t *testing.T) {
	expectedMovie := &model.Movie{MovieId: "tt0111161", Name: "The Shawshank Redemption", Duration: "2h22m", Plot: "Two imprisoned men bond."}

	tests := []struct {
		name      string
		setupMock func(gateway *moviemocks.MockMovieGateWay)
		wantMovie *model.Movie
		wantErr   bool
	}{
		{
			name: "Should return movie when gateway succeeds",
			setupMock: func(gateway *moviemocks.MockMovieGateWay) {
				gateway.On("MovieById", mock.Anything, mock.AnythingOfType("string")).
					Return(expectedMovie, nil).Once()
			},
			wantMovie: expectedMovie,
		},
		{
			name: "Should propagate error when movie gateway fails",
			setupMock: func(gateway *moviemocks.MockMovieGateWay) {
				gateway.On("MovieById", mock.Anything, mock.AnythingOfType("string")).
					Return(&model.Movie{}, fmt.Errorf("gateway error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			showRepo := servicemocks.NewMockShowRepository(t)
			gateway := moviemocks.NewMockMovieGateWay(t)
			tc.setupMock(gateway)

			svc := NewShowService(showRepo, gateway)
			got, err := svc.GetMovieById(context.Background(), "")

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantMovie, got)
			}
		})
	}
}
