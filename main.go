package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies) //convert all the movies in JSON format and return
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)             //type of parse for the received data
	for index, item := range movies { //for index and item in range of movies (similar to python)
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //append n+1 to end -> 1 to n thus skipping the current value
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies { //since we dont need a index here we have to pass nothing i.e. _ so that go doesnt gives an error
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	w.Header().Set("Conten-type", "application/json")
	//params
	params := mux.Vars(r)
	//loopover the movies range
	for index, item := range movies {
		if item.ID == params["id"] {
			//delete the movies with the i.d. that you've sent
			movies = append(movies[:index], movies[index+1:]...)
			//add a new movie: the movie sent through  posytman
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

func main() {

	r := mux.NewRouter()

	//CREATING SOME MOVIES SO IT IS NOT INITIALLY NULL
	movies = append(movies, Movie{ID: "1", Isbn: "45623", Title: "Movie1", Director: &Director{FirstName: "John", LastName: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "54585", Title: "Movie2", Director: &Director{FirstName: "Dohn", LastName: "Joe"}})

	//DECLARING ALL THE POSSIBLE ROUTES FO THE SERVER.
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Server startred on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r)) //CREATE AND START A SERVER AT 'localhost:8000'
}
