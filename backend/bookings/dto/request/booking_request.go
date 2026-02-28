package request

import (
	"skyfox/bookings/model"
)

type BookingRequest struct {
	Date      string         `json:"date"  binding:"required,datetime=2006-01-02"`
	ShowId    int            `json:"showId" binding:"required,gte=1"`
	Customer  model.Customer `json:"customer" binding:"required"`
	NoOfSeats int            `json:"noOfSeats" binding:"maxSeats"`
}
