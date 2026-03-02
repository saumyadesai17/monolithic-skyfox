package controller

import (
	"context"
	"net/http"
	"skyfox/common/logger"
	ae "skyfox/error"

	"github.com/gin-gonic/gin"
)

type RevenueService interface {
	RevenueOn(context.Context, string) (float64, error)
}

type RevenueController struct {
	revenueService RevenueService
}

func NewRevenueController(revenueService RevenueService) *RevenueController {
	return &RevenueController{
		revenueService: revenueService,
	}
}

// Revenue godoc
//
//		@Summary		Revenue
//		@Description	get revenue by date
//		@Tags			Revenue
//	 @security	BasicAuth
//		@Accept			json
//		@Produce		json
//	 @param Authorization header string true "Enter basic auth"
//	 @Param  date  query string true "to get shows"
//		@Success		200	{number}	string
//		@Failure		400	{object}	ae.AppError
//		@Failure		404	{object}	ae.AppError
//		@Failure		500	{object}	ae.AppError
//		@Router			/revenue [get]
func (rh *RevenueController) GetRevenue(c *gin.Context) {
	date := c.Request.URL.Query().Get("date")
	revenue, responseError := rh.revenueService.RevenueOn(c.Request.Context(), date)

	if responseError != nil {
		err := responseError.(*ae.AppError)
		logger.Error("%s", err.UnWrap().Error())
		c.AbortWithStatusJSON(err.HTTPCode(), err)
	}

	c.IndentedJSON(http.StatusOK, revenue)
}
