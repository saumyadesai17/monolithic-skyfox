package request

import (
	"skyfox/bookings/constants"
	"skyfox/bookings/model"
	"skyfox/common/middleware/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookingRequest(t *testing.T) {
	t.Run("should provide valid date", func(t *testing.T) {

		bookingRequest := BookingRequest{
			Date:   "",
			ShowId: 1,
			Customer: model.Customer{
				Name:        "john",
				PhoneNumber: "9999988888",
			},
			NoOfSeats: 2,
		}

		validate := new(validator.DtoValidator)
		err := validate.ValidateStruct(bookingRequest)

		assert.NotNil(t, err)
	})
	t.Run("should not allow to book more seats than allowed", func(t *testing.T) {

		maxSeats := constants.MAX_NO_OF_SEATS_PER_BOOKING
		bookingRequest := BookingRequest{
			Date:   "2022-11-20",
			ShowId: 1,
			Customer: model.Customer{
				Name:        "john",
				PhoneNumber: "9999988888",
			},
			NoOfSeats: maxSeats + 1,
		}

		validate := new(validator.DtoValidator)
		err := validate.ValidateStruct(bookingRequest)

		assert.NotNil(t, err)
	})
}
