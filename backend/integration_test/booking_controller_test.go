package integrationtest

import (
	"context"
	"net/http"
	"skyfox/bookings/constants"
	"skyfox/bookings/controller"
	req "skyfox/bookings/dto/request"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/bookings/repository/testdata"
	"skyfox/bookings/service"
	"skyfox/common/middleware/security"
	db "skyfox/integration_test/db"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var bookingsPath = constants.BookingEndPoint

// TODO - standardize the test names
func Test_WhenCreateBooking_ItShouldReturnBookingConfirmationResponse(t *testing.T) {

	tearDown := bookingControllerTestSetup(t)
	defer tearDown(t)

	db := db.GetDB()
	userService := service.NewUserService(repository.NewUserRepository(db))
	bookingService := service.NewBookingService(repository.NewBookingRepository(db), repository.NewShowRepository(db))
	bookingService.SetCustomerRepository(repository.NewCustomerRepository(db))
	handler := controller.NewBookingController(bookingService)

	engine, request := getEngine()
	engine.POST(bookingsPath, security.Authenticate(userService), handler.CreateBooking)

	payload := req.BookingRequest{
		Date:      "2022-10-13",
		ShowId:    1,
		Customer:  model.Customer{Name: "John", PhoneNumber: "6543276543"},
		NoOfSeats: 2,
	}
	bookingResponse := `{
		"id": 1,
		"customerName": "John",
		"showDate": "2022-10-13",
		"startTime": "18:00:00",
		"amountPaid": 600,
		"noOfSeats": 2
	}`

	request.POST(bookingsPath).SetDebug(true).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
		SetJSONInterface(payload).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
			require.JSONEq(t, bookingResponse, r.Body.String())
		})
}

func Test_WhenCreateBooking_ForSeatsMoreThanAllowed_ItShouldNotBook(t *testing.T) {
	tearDown := bookingControllerTestSetup(t)
	defer tearDown(t)

	db := db.GetDB()
	userService := service.NewUserService(repository.NewUserRepository(db))
	bookingService := service.NewBookingService(repository.NewBookingRepository(db), repository.NewShowRepository(db))
	bookingService.SetCustomerRepository(repository.NewCustomerRepository(db))
	handler := controller.NewBookingController(bookingService)

	engine, request := getEngine()
	engine.POST(bookingsPath, security.Authenticate(userService), handler.CreateBooking)

	maxSeats := constants.MAX_NO_OF_SEATS_PER_BOOKING

	payload := req.BookingRequest{
		Date:      "2022-10-13",
		ShowId:    1,
		Customer:  model.Customer{Name: "John", PhoneNumber: "6543276543"},
		NoOfSeats: maxSeats + 1,
	}

	request.POST(bookingsPath).SetDebug(true).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
		SetJSONInterface(payload).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}

func Test_WhenCreateBooking_And_MaxCapacityOfSeatsExceeds_ItShouldNotBook(t *testing.T) {
	tearDown := bookingControllerTestSetup(t)
	defer tearDown(t)

	db := db.GetDB()
	userService := service.NewUserService(repository.NewUserRepository(db))
	bookingService := service.NewBookingService(repository.NewBookingRepository(db), repository.NewShowRepository(db))
	bookingService.SetCustomerRepository(repository.NewCustomerRepository(db))
	handler := controller.NewBookingController(bookingService)

	engine, request := getEngine()
	engine.POST(bookingsPath, security.Authenticate(userService), handler.CreateBooking)

	setupBookingSeatsForSameShow(engine, request)

	payload := req.BookingRequest{
		Date:      "2022-10-13",
		ShowId:    1,
		Customer:  model.Customer{Name: "John", PhoneNumber: "6543276543"},
		NoOfSeats: 11,
	}

	request.POST(bookingsPath).SetDebug(true).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
		SetJSONInterface(payload).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}

func bookingControllerTestSetup(t *testing.T) func(*testing.T) {
	db := db.GetDB()
	gormDB := db.GormDB()

	// Clean up existing data first to avoid conflicts
	gormDB.Exec("DELETE FROM BOOKING")
	gormDB.Exec("DELETE FROM SHOW")
	gormDB.Exec("DELETE FROM SLOT")
	gormDB.Exec("DELETE FROM USERTABLE")

	// Reset sequences to ensure IDs start from 1
	gormDB.Exec("ALTER SEQUENCE show_id_seq RESTART WITH 1")
	gormDB.Exec("ALTER SEQUENCE booking_id_seq RESTART WITH 1")

	// Now create test data
	gormDB.Create(testdata.Shows)
	user := model.NewUser("admin", "password")
	userRepo := repository.NewUserRepository(db)
	err := userRepo.Create(context.Background(), &user)

	if err != nil {
		panic(err)
	}
	return func(t *testing.T) {
		gormDB.Exec("DELETE FROM BOOKING")
		gormDB.Exec("DELETE FROM SHOW")
		gormDB.Exec("DELETE FROM SLOT")
		gormDB.Exec("DELETE FROM USERTABLE")
	}
}

func setupBookingSeatsForSameShow(e *gin.Engine, r *gofight.RequestConfig) {
	payload := req.BookingRequest{
		Date:      "2022-10-13",
		ShowId:    1,
		Customer:  model.Customer{Name: "John", PhoneNumber: "6543276543"},
		NoOfSeats: constants.MAX_NO_OF_SEATS_PER_BOOKING,
	}

	for i := 0; i < 6; i++ {
		r.POST(bookingsPath).SetDebug(true).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
			SetJSONInterface(payload).
			Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			})
	}
}
