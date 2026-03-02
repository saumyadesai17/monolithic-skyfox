package controller

import (
	"context"
	"net/http"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/dto/response"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	"skyfox/common/middleware/validator"
	ae "skyfox/error"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Signup(ctx context.Context, req request.SignupRequest) (*model.UserAccount, error)
}

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Signup godoc
//
//	@Summary		Customer Signup
//	@Description	Register a new customer account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.SignupRequest	true	"Signup payload"
//	@Success		201		{object}	response.SignupResponse
//	@Failure		400		{object}	ae.AppError
//	@Failure		409		{object}	ae.AppError
//	@Failure		500		{object}	ae.AppError
//	@Router			/api/v1/auth/signup [post]
func (ac *AuthController) Signup(c *gin.Context) {
	var req request.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("signup validation failed: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, validator.HandleStructValidationError(err))
		return
	}

	user, err := ac.authService.Signup(c.Request.Context(), req)
	if err != nil {
		appErr := err.(*ae.AppError)
		logger.Error("signup error: %v", appErr)
		c.AbortWithStatusJSON(appErr.HTTPCode(), appErr)
		return
	}

	c.IndentedJSON(http.StatusCreated, response.NewSignupResponse(user.Id, user.Name, user.Phone))
}
