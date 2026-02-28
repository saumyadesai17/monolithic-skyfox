package movieservice

import (
	"context"
	"skyfox/_mocks/repomocks"
	"skyfox/bookings/model"
	"skyfox/config"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/h2non/gock.v1"
)

func Test_ReturnsMovie_When_MovieServiceIsInvoked(t *testing.T) {
	want := &model.Movie{
		MovieId: "movie_id",
		Name: "movie_name",
		Duration: "1h30m0s",
		Plot: "movie plot in short",
	}
	
	movieGatewayRepo := repomocks.MovieGateWay{}
	movieGatewayRepo.On("MovieById", mock.AnythingOfType("string")).Return(want, nil).Once()
	defer gock.Off()
    body := `{"Title":"movie_name","Year":"2018","Rated":"PG-13","Released":"06 Apr 2018","Runtime":"90 min","Genre":"Drama, Horror, Sci-Fi","Director":"John Krasinski","Writer":"Bryan Woods (screenplay by), Scott Beck (screenplay by), John Krasinski (screenplay by), Bryan Woods (story by), Scott Beck (story by)","Actors":"Emily Blunt, John Krasinski, Millicent Simmonds, Noah Jupe","Plot":"movie plot in short","Language":"English, American Sign Language","Country":"USA","Awards":"Nominated for 1 Oscar. Another 34 wins & 108 nominations.","Poster":"https://m.media-amazon.com/images/M/MV5BMjI0MDMzNTQ0M15BMl5BanBnXkFtZTgwMTM5NzM3NDM@._V1_SX300.jpg","Ratings":[{"Source":"Internet Movie Database","Value":"7.5/10"},{"Source":"Rotten Tomatoes","Value":"95%"},{"Source":"Metacritic","Value":"82/100"}],"Metascore":"82","imdbRating":"7.5","imdbVotes":"379,472","imdbID":"movie_id","Type":"movie","DVD":"N/A","BoxOffice":"N/A","Production":"N/A","Website":"N/A","Response":"True"}`
    gock.New("http://localhost:4567").Get("/movies/*").Persist().Reply(200).BodyString(body)

	movieGateway := NewMovieGateway(config.MovieGatewayConfig{MovieServiceHost: "http://localhost:4567/"})

	got,err := movieGateway.MovieById(context.Background(),"")

	assert.Nil(t,err)
	assert.Equal(t, want, got)
}
