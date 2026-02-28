package movieservice

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"skyfox/bookings/model"
	"skyfox/config"

	ae "skyfox/error"

	"github.com/bborbe/http/requestbuilder"
)

type MovieGateWay interface {
	MovieById(ctx context.Context, id string) (*model.Movie, error)
}

type movieGateway struct {
	config config.MovieGatewayConfig
}

func NewMovieGateway(cfg config.MovieGatewayConfig) MovieGateWay {
	return &movieGateway{
		config: config.MovieGatewayConfig{
			MovieServiceHost: cfg.MovieServiceHost,
		},
	}
}

func (d *movieGateway) MovieById(ctx context.Context, id string) (*model.Movie, error) {

	var http http.Client
	var err error

	request, err := requestbuilder.NewHTTPRequestBuilder(d.config.MovieServiceHost + "movies/" + id).Build()
	if err != nil {
		return &model.Movie{}, ae.InternalServerError("InternalServerError", "could not parse movie service url", err)
	}

	httpResponse, err := http.Do(request)
	if err != nil {
		return &model.Movie{}, ae.InternalServerError("InternalServerError", "could not retrieve the movie detail", err)
	}
	defer httpResponse.Body.Close()
	responseBody, _ := ioutil.ReadAll(httpResponse.Body)

	var movieResponse MovieServiceResponse
	err = json.Unmarshal(responseBody, &movieResponse)

	if err != nil {
		return &model.Movie{}, ae.InternalServerError("InternalServerError", "failed to parse the movie details", err)
	}

	movie, err := movieResponse.ToMovie()
	if err != nil {
		return &model.Movie{}, err
	}
	return movie, nil
}
