package integrationtest

import (

	// req "skyfox/bookings/model/request"

	"context"
	"net/http"
	"skyfox/bookings/constants"
	"skyfox/bookings/controller"
	"skyfox/bookings/model"
	"skyfox/bookings/repository"
	"skyfox/bookings/service"
	"skyfox/common/middleware/security"
	db "skyfox/integration_test/db"
	"testing"

	"github.com/appleboy/gofight/v2"
	"gotest.tools/assert"
)

var revenuePath = constants.RevenueEndPoint

func Test_WhenGetRevenue_ItShouldReturnRevenue(t *testing.T) {
	tearDown := revenueControllerTestSetup(t)
	defer tearDown(t)

	db := db.GetDB()
	gormDB := db.GormDB()
	gormDB.Create(dummyShows)

	showRepo := repository.NewShowRepository(db)
	showOne, _ := showRepo.FindById(context.Background(), 1)
	showTwo, _ := showRepo.FindById(context.Background(), 2)
	showThree, _ := showRepo.FindById(context.Background(), 3)

	testCustomer := model.Customer{
		Name:        "john",
		PhoneNumber: "9999988888",
	}

	bookingRepo := repository.NewBookingRepository(db)
	newBooking := model.NewBooking(showOne.Date, showOne, testCustomer, 2, 400.00)
	bookingRepo.Create(context.Background(), &newBooking)

	newBooking = model.NewBooking(showTwo.Date, showTwo, testCustomer, 3, 450.00)
	bookingRepo.Create(context.Background(), &newBooking)

	newBooking = model.NewBooking(showThree.Date, showThree, testCustomer, 1, 250.00)
	bookingRepo.Create(context.Background(), &newBooking)

	handler := controller.NewRevenueController(service.NewRevenueService(repository.NewBookingRepository(db), repository.NewShowRepository(db)))
	userService := service.NewUserService(repository.NewUserRepository(db))

	engine, request := getEngine()
	engine.GET(revenuePath, security.Authenticate(userService), handler.GetRevenue)

	request.GET(revenuePath).
		SetQuery(gofight.H{"date": "2022-10-13"}).
		SetDebug(true).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "850", r.Body.String())
		})
}

func Test_WhenGetRevenue_ForNoBooking_ItShouldReturnZero(t *testing.T) {

	tearDown := revenueControllerTestSetup(t)
	defer tearDown(t)

	db := db.GetDB()
	userService := service.NewUserService(repository.NewUserRepository(db))

	handler := controller.NewRevenueController(service.NewRevenueService(repository.NewBookingRepository(db), repository.NewShowRepository(db)))

	engine, request := getEngine()
	engine.GET(revenuePath, security.Authenticate(userService), handler.GetRevenue)

	request.GET(revenuePath).SetDebug(true).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "0", r.Body.String())
		})
}

func revenueControllerTestSetup(t *testing.T) func(*testing.T) {
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
	gormDB.Create(dummyShows)
	user := model.NewUser("admin", "password")
	userRepo := repository.NewUserRepository(db)
	userRepo.Create(context.Background(), &user)

	return func(t *testing.T) {
		gormDB.Exec("DELETE FROM BOOKING")
		gormDB.Exec("DELETE FROM SHOW")
		gormDB.Exec("DELETE FROM SLOT")
		gormDB.Exec("DELETE FROM USERTABLE")
	}
}

var dummyShows = []model.Show{
	{
		Id:      1,
		MovieId: "movie_1",
		Date:    "2022-10-13",
		SlotId:  1,
		Slot: model.Slot{
			Id:        1,
			Name:      "slot_1",
			StartTime: "18:00:00",
			EndTime:   "21:30:00",
		},
		Cost: 200.00,
	},
	{
		Id:      2,
		MovieId: "movie_2",
		Date:    "2022-10-13",
		SlotId:  2,
		Slot: model.Slot{
			Id:        2,
			Name:      "slot_2",
			StartTime: "22:30:00",
			EndTime:   "02:00:00",
		},
		Cost: 150.00,
	},
	{
		Id:      3,
		MovieId: "movie_3",
		Date:    "2022-10-14",
		SlotId:  2,
		Slot: model.Slot{
			Id:        2,
			Name:      "slot_2",
			StartTime: "22:30:00",
			EndTime:   "02:00:00",
		},
		Cost: 250.00,
	},
}
