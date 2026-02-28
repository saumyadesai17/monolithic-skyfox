package service

import (
	"context"
)

type revenueService struct {
	showRepo ShowRepository
	bookRepo BookingRepository
}

func NewRevenueService(bookingRepository BookingRepository, showRepository ShowRepository) *revenueService {
	return &revenueService{
		showRepo: showRepository,
		bookRepo: bookingRepository,
	}
}

func (rs *revenueService) RevenueOn(ctx context.Context, date string) (float64, error) {
	shows, err := rs.showRepo.GetAllShowsOn(ctx, date)
	if err != nil {
		return 0, err
	}

	if len(shows) == 0 {
		return 0, nil
	}

	var showIds []int
	for _, show := range shows {
		showIds = append(showIds, show.Id)
	}

	revenue := rs.bookRepo.BookingAmountByShows(ctx, showIds)
	return revenue, nil
}
