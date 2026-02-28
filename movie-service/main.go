
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Movie represents a movie structure
type Movie struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	ImdbRating string   `json:"imdbRating"`
	ImdbVotes  string   `json:"imdbVotes"`
	ImdbID     string   `json:"imdbID"`
	Type       string   `json:"Type"`
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"`
}

// Rating represents a movie rating
type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

var movies []Movie

// loadMovies loads movies from JSON file
func loadMovies() error {
	data, err := ioutil.ReadFile("movies.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &movies)
}

// getMovies returns all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// getMovieByID returns a single movie by ID
func getMovieByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, movie := range movies {
		if movie.ImdbID == id {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Movie with requested ID not found"})
}

// getRoot returns a random quote with color
func getRoot(w http.ResponseWriter, r *http.Request) {
	quotes := []string{
		"The only way to do great work is to love what you do - Steve Jobs",
		"Innovation distinguishes between a leader and a follower - Steve Jobs",
		"Stay hungry, stay foolish - Steve Jobs",
		"Life is what happens when you're busy making other plans - John Lennon",
		"The future belongs to those who believe in the beauty of their dreams - Eleanor Roosevelt",
	}

	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A", "#98D8C8", "#F7DC6F", "#BB8FCE"}

	rand.Seed(time.Now().UnixNano())
	quote := quotes[rand.Intn(len(quotes))]
	color := colors[rand.Intn(len(colors))]

	html := fmt.Sprintf(`
		<html>
		<body style="display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; font-family: Arial, sans-serif;">
			<p style="text-align: center; font-size: 2em; color: %s; max-width: 80%%;">
				%s
			</p>
		</body>
		</html>
	`, color, quote)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

// notFound handles 404 errors
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "There is nothing to do here.! 404!")
}

func main() {
	// Load movies from JSON file
	if err := loadMovies(); err != nil {
		log.Fatal("Error loading movies:", err)
	}

	log.Printf("Loaded %d movies from movies.json\n", len(movies))

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/", getRoot).Methods("GET")
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(notFound)

	// Start server
	port := "4567"
	log.Printf("Movie service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}


