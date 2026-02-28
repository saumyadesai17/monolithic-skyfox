package movieservice

import (
	"fmt"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	ae "skyfox/error"
	"strings"
	"time"
)

type MovieServiceResponse struct {
	ImdbId  string `json:"imdbid"`
	Title   string `json:"title"`
	RunTime string `json:"runtime"`
	Plot    string `json:"plot"`
}

func (m MovieServiceResponse) ToMovie() (*model.Movie, error) {
	runtime := strings.Split(m.RunTime, " ")[0]
	duration, err := time.ParseDuration(runtime + "m")

	if err != nil {
		logger.Error(fmt.Sprintf("failed to get the run time of the movie %s", m.Title), err)
		return &model.Movie{}, ae.UnProcessableError("MovieCreationFailed", "Movie creation failed due to unknown reason", err)
	}
	movie := model.NewMovie(m.ImdbId, m.Title, duration.String(), m.Plot)
	return movie, nil
}
