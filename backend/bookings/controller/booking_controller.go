package controller

import (
	"context"
	"math"
	"net/http"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/dto/response"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/common/logger"
	"skyfox/common/middleware/validator"

	ae "skyfox/error"

	"github.com/gin-gonic/gin"
)

type BookingService interface {
	Book(context.Context, request.BookingRequest) (*model.Booking, error)
	SetCustomerRepository(repository.CustomerRepository)
}

type BookingController struct {
	bookingService BookingService
}

func NewBookingController(bookingService BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

// booking godoc
//
//		@Summary		Booking
//		@Description	Book a Ticket
//		@Tags			Booking
//		@security	BasicAuth
//	 @param Authorization header string true "Enter basic auth"
//		@Accept			json
//		@Produce		json
//		@Param			request	body		request.BookingRequest	true	"Book Ticket"
//		@Success		200	{object}	response.BookingConfirmationResponse
//		@Failure		400	{object}	ae.AppError
//		@Failure		404	{object}	ae.AppError
//		@Failure		500	{object}	ae.AppError
//		@Router			/bookings [POST]
func (bh *BookingController) CreateBooking(c *gin.Context) {
	var newBookingRequest request.BookingRequest

	if err := c.ShouldBindJSON(&newBookingRequest); err != nil {
		logger.Error("bad request %v", validator.HandleStructValidationError(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, validator.HandleStructValidationError(err))
		return
	}

	bookingResponse, responseError := bh.bookingService.Book(c.Request.Context(), newBookingRequest)

	if responseError != nil {
		err := responseError.(*ae.AppError)
		logger.Error("error occurred. %v", err)
		c.AbortWithStatusJSON(err.HTTPCode(), err)
		return
	}

	bookingConfirmation := response.NewBookingConfirmationResponse(
		bookingResponse.Id,
		bookingResponse.Customer.Name,
		bookingResponse.Show.Date,
		bookingResponse.Show.Slot.StartTime,
		math.Round(bookingResponse.AmountPaid*100)/100,
		bookingResponse.NoOfSeats,
	)
	c.IndentedJSON(http.StatusCreated, bookingConfirmation)
}
