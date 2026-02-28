package controller

import (
	"context"
	"net/http"
	"skyfox/bookings/model"
	ae "skyfox/error"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	UserDetails(context.Context, string) (model.User, error)
}

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// login godoc
//
//		@Summary		Login
//		@Description	login
//		@Tags			login
//		@Accept			json
//		@Produce		json
//	 @param Authorization header string true "Enter basic auth"
//		@Success		200	{string}	string
//		@Failure		401	{object}	ae.AppError
//		@Failure		404	{object}	ae.AppError
//		@Failure		500	{object}	ae.AppError
//		@Router			/login [get]
func (uh *UserController) Login(c *gin.Context) {

	username, _, _ := c.Request.BasicAuth()

	user, err := uh.userService.UserDetails(c.Request.Context(), username)
	if err != nil {
		appError := err.(*ae.AppError)
		c.AbortWithStatusJSON(appError.HTTPCode(), appError)
		return
	}

	c.JSON(http.StatusOK, user.Username)
}
