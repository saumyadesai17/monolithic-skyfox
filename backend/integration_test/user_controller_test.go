package integrationtest

import (
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
	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

var loginPath = constants.LoginEndPoint

func Test_WhenLogin_ItShouldReturnOK(t *testing.T) {

	engine, request, tearDown := userControllerTestSetUp(t)
	defer tearDown(t)

	request.GET(loginPath).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc3dvcmQ="}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func Test_WhenLoginWithInvalidCredential_ItShouldReturnUnauthorized(t *testing.T) {

	engine, request, tearDown := userControllerTestSetUp(t)
	defer tearDown(t)

	request.GET(loginPath).SetHeader(gofight.H{"Authorization": "Basic YWRtaW46cGFzc4evcmQ="}).
		Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusUnauthorized, r.Code)
		})
}

func userControllerTestSetUp(t *testing.T) (*gin.Engine, *gofight.RequestConfig, func(*testing.T)) {
	db := db.GetDB()
	gormDB := db.GormDB()

	// Clean up existing data first to avoid conflicts
	gormDB.Exec("DELETE FROM USERTABLE")

	// Reset sequence to ensure IDs start from 1
	gormDB.Exec("ALTER SEQUENCE usertable_id_seq RESTART WITH 1")

	// Now create test data
	user := model.NewUser("admin", "password")
	userRepo := repository.NewUserRepository(db)
	userRepo.Create(context.Background(), &user)
	userService := service.NewUserService(repository.NewUserRepository(db))

	handler := controller.NewUserController(service.NewUserService(repository.NewUserRepository(db)))

	engine, request := getEngine()
	engine.GET(loginPath, security.Authenticate(userService), handler.Login)

	return engine, request, func(t *testing.T) {
		gormDB.Exec("DELETE FROM USERTABLE")
	}
}
