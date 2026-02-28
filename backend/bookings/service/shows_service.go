package service

import (
	"context"
	"skyfox/bookings/model"
	movieservice "skyfox/movieservice/movie_gateway"
)

type ShowRepository interface {
	GetAllShowsOn(ctx context.Context, date string) ([]model.Show, error)
	FindById(ctx context.Context, id int) (model.Show, error)
}

type showService struct {
	showRepo     ShowRepository
	movieGateway movieservice.MovieGateWay
}

func NewShowService(showRepository ShowRepository, gateway movieservice.MovieGateWay) *showService {
	return &showService{
		showRepo:     showRepository,
		movieGateway: gateway,
	}
}

func (s *showService) GetShows(ctx context.Context, date string) ([]model.Show, error) {
	shows, err := s.showRepo.GetAllShowsOn(ctx, date)
	if err != nil {
		return nil, err
	}
	return shows, nil
}

func (s *showService) GetMovieById(ctx context.Context, movieId string) (*model.Movie, error) {
	movie, err := s.movieGateway.MovieById(ctx, movieId)
	if err != nil {
		return &model.Movie{}, err
	}
	return movie, nil
}
