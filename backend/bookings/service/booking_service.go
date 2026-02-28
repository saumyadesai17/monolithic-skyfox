package service

import (
	"context"
	"errors"
	"fmt"
	"skyfox/bookings/constants"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/common/logger"
	ae "skyfox/error"
)

type BookingRepository interface {
	Create(context.Context, *model.Booking) error
	BookedSeatsByShow(context.Context, int) int
	BookingAmountByShows(context.Context, []int) float64
}

type bookingService struct {
	bookRepo     BookingRepository
	showRepo     ShowRepository
	customerRepo repository.CustomerRepository
}

func NewBookingService(bookingRepository BookingRepository, showRepository ShowRepository) *bookingService {
	return &bookingService{
		bookRepo:     bookingRepository,
		showRepo:     showRepository,
		customerRepo: nil, // Will be set by SetCustomerRepository
	}
}

func (b *bookingService) SetCustomerRepository(customerRepository repository.CustomerRepository) {
	b.customerRepo = customerRepository
}

func (b *bookingService) Book(ctx context.Context, newBookingRequest request.BookingRequest) (*model.Booking, error) {
	show, err := b.showRepo.FindById(ctx, newBookingRequest.ShowId)
	if err != nil {
		logger.Error("error occurred. %v", err)
		return nil, err
	}

	seatsRequested := newBookingRequest.NoOfSeats
	if err != nil {
		logger.Error("error occurred. %v", fmt.Errorf("can't process no. of seats: %v", err))
		return nil, ae.BadRequestError("BadRequest", "no. of seats requested is invalid", fmt.Errorf("can't process no. of seats: %v", err))
	}

	maxBookingAllowed := constants.MAX_NO_OF_SEATS_PER_BOOKING
	if seatsRequested > maxBookingAllowed {
		logger.Error("error occurred. %v", fmt.Errorf("can't book seats: %v", errors.New("maximum seats per booking exceeded")))
		return nil, ae.BadRequestError("BadRequest", fmt.Sprintf("maximum seats per booking is %d", constants.MAX_NO_OF_SEATS_PER_BOOKING),
			fmt.Errorf("can't book seats: %v", errors.New("maximum seats per booking exceeded")))
	}

	amountPaid := show.Cost * float64(seatsRequested)

	availableSeats := b.availableSeats(ctx, show.Id)
	if availableSeats < seatsRequested {
		logger.Error("error occurred. %v", fmt.Sprintf("seats requested %d, seats available %d", seatsRequested, availableSeats))
		return nil, ae.BadRequestError("NoSeatAvailable", fmt.Sprintf("seats requested %d, seats available %d", seatsRequested, availableSeats),
			fmt.Errorf("seats available %d", availableSeats))
	}
	// First save the customer
	customer := newBookingRequest.Customer
	if b.customerRepo != nil {
		err = b.customerRepo.Create(ctx, &customer)
		if err != nil {
			logger.Error("error occurred while creating customer. %v", err)
			// Continue even if customer already exists
		}
	}

	newBooking := model.NewBooking(newBookingRequest.Date, show, customer, seatsRequested, amountPaid)
	err = b.bookRepo.Create(ctx, &newBooking)
	if err != nil {
		logger.Error("error occurred. %v", err)
		return nil, err
	}
	return &newBooking, nil
}

func (b *bookingService) availableSeats(ctx context.Context, showId int) int {
	bookedSeats := b.bookRepo.BookedSeatsByShow(ctx, showId)
	return constants.TOTAL_NO_OF_SEATS - bookedSeats
}
