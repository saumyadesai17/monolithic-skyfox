package security

import (
	"fmt"
	"net/http"
	"skyfox/bookings/controller"

	"skyfox/bookings/model"
	"skyfox/common/logger"
	ae "skyfox/error"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate(userService controller.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			err := ae.BadRequestError("BadRequest", "only basic auth is supported", fmt.Errorf("only basic auth is supported for sign in"))
			logger.Error(err.UnWrap().Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		if isParameterEmpty(username, password) {
			err := ae.BadRequestError("BadRequest", "username and password missing", fmt.Errorf("username and password can't be empty"))
			logger.Error(err.UnWrap().Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		user, err := userService.UserDetails(c, username)
		if err != nil {
			err := err.(*ae.AppError)
			logger.Error(err.UnWrap().Error())
			c.AbortWithStatusJSON(err.HTTPCode(), err)
			return
		}

		if isNotValid(password, user) {
			err := ae.InvalidCredentialsError("WrongCredentials", "username or password is wrong", fmt.Errorf("credentials do not match, failed to authenticate"))
			logger.Error(err.UnWrap().Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		c.Next()
	}
}

func isParameterEmpty(username string, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func isNotValid(password string, user model.User) bool {
	return user == model.User{} || !(user.Password == password)
}
