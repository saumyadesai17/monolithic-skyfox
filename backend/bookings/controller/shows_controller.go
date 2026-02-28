package controller

import (
	"context"

	"net/http"
	"skyfox/bookings/dto/response"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	ae "skyfox/error"

	"github.com/gin-gonic/gin"
)

type ShowService interface {
	GetShows(context.Context, string) ([]model.Show, error)
	GetMovieById(context.Context, string) (*model.Movie, error)
}

type showController struct {
	showService ShowService
}

func NewShowController(showService ShowService) *showController {
	return &showController{
		showService: showService,
	}
}

// Shows godoc
//
//		@Summary		Shows
//		@Description	get shows by date
//		@Tags			Shows
//		@Accept			json
//		@Produce		json
//	 @security	BasicAuth
//	 @param Authorization header string true "Enter basic auth"
//	 @Param  date  query string true "to get shows"
//		@Success		200	{object}	response.ShowResponse
//		@Failure		400	{object}	ae.AppError
//		@Failure		404	{object}	ae.AppError
//		@Failure		500	{object}	ae.AppError
//		@Router			/shows [get]
func (sh *showController) Shows(c *gin.Context) {
	date := c.Request.URL.Query().Get("date")

	shows, responseError := sh.showService.GetShows(c.Request.Context(), date)
	if responseError != nil {
		err := responseError.(*ae.AppError)
		logger.Error(err.UnWrap().Error())
		c.AbortWithStatusJSON(err.HTTPCode(), err)
	}

	var showResponses []response.ShowResponse
	for _, show := range shows {
		show_response, responseError := sh.constructShowResponse(c.Request.Context(), show)
		if responseError != nil {
			err := responseError.(*ae.AppError)
			// logger.Error(err.UnWrap().Error())
			c.AbortWithStatusJSON(err.HTTPCode(), err)
			return
		}
		showResponses = append(showResponses, *show_response)
	}
	c.IndentedJSON(http.StatusOK, showResponses)
}

func (sh *showController) constructShowResponse(ctx context.Context, s model.Show) (*response.ShowResponse, error) {
	movie, err := sh.showService.GetMovieById(ctx, s.MovieId)
	if err != nil {
		return &response.ShowResponse{}, err
	}
	return response.NewShowResponse(*movie, s.Slot, s), nil
}
