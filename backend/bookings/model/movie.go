package model

type Movie struct {
	MovieId  string `json:"id"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
	Plot     string `json:"plot"`
}

func NewMovie(id string, name string, duration string, plot string) *Movie {
	return &Movie{
		MovieId:  id,
		Name:     name,
		Duration: duration,
		Plot:     plot,
	}
}
