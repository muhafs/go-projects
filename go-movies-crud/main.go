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
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var movies = []Movie{
	{ID: "1", ISBN: "12345", Title: "Mov1", Director: &Director{FirstName: "John", LastName: "Doe"}},
	{ID: "2", ISBN: "32345", Title: "Mov2", Director: &Director{FirstName: "Steve", LastName: "Smith"}},
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(movies) // convert Movies array into JSON then return it
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie) // convert Movie object into JSON then return it
			return
		}
	}

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	movie := Movie{}
	json.NewDecoder(r.Body).Decode(&movie) // convert JSON request of Movie struct into Movie object then assign it
	//? Gin = c.BindJSON(&movie)

	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movies) // convert Movies array into JSON then return it
	//? Gin = c.IndentedJSON(200, movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, movie := range movies {
		if movie.ID == params["id"] {
			// delete movie
			movies = append(movies[:i], movies[i+1:]...)

			// append new movie
			newMovie := Movie{}
			json.NewDecoder(r.Body).Decode(&newMovie) // convert JSON request of Movie struct into Movie object then assign it
			newMovie.ID = params["id"]
			movies = append(movies, newMovie)

			// return
			json.NewEncoder(w).Encode(movies) // convert Movies array into JSON then return it
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)

			json.NewEncoder(w).Encode(movies) // convert Movies array into JSON then return it
			return
		}
	}
}

func main() {
	/*
		* Define router
		- Example:
		-- Gin = gin.Default()
		-- Chi = chi.NewRouter()
		-- Mux = mux.NewRouter()
	*/
	r := mux.NewRouter()

	/*
		* Movies Routes
		- Example:
		-- Gin = r.GET("/", someHandler)
		-- Chi = r.Get("/", someHandler)
		-- Mux = r.HandleFunc("/", someHandler).Methods("GET")
	*/
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	//* Start server
	port := ":8000"
	fmt.Printf("Server is running at http://127.0.0.1%s\n", port)

	err := http.ListenAndServe(port, r) //? Gin = r.Run(":8080")
	if err != nil {
		log.Fatal("Couldn't start server:", err)
	}
}
