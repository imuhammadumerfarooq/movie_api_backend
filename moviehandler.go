package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Movie struct definition
type Movie struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
	Rating string `json:"rating"`
	Image  string `json:"coverImage"`
	Price  int    `json:"price"`
}

// Create a response struct to include movies and the total count
type MovieResponse struct {
	Movies []Movie `json:"movies"`
	Count  int     `json:"count"`
}

// Function to count all movies in the database
func countMovies() (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM movies")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Handler function to get all movies
func getMovies(c echo.Context) error {

	// Get the limit from query parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 8
	}

	// Get the offset from query parameters
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// Read all movies form database with limit and offset
	movies, err := readMovies(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Get the total count of movies from the database
	count, err := countMovies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Return the list of movies and the total count with a status OK
	return c.JSON(http.StatusOK, MovieResponse{
		Movies: movies,
		Count:  count,
	})

}

// Handler function to get a single movie by ID
func getMovie(c echo.Context) error {

	id := c.Param("id")

	// Read movie by specific ID from database
	movie, err := readMovieById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, "Movie not found")
	}

	return c.JSON(http.StatusOK, movie)
}

// Handler function to create a new movie
func createMovie(c echo.Context) error {
	var movie Movie

	// Bind the incoming JSON request to the Movie Struct
	if err := c.Bind(&movie); err != nil {
		log.Print("Stuck here")
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Write the new movie to the database
	if err := writeMovie(movie); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, movie)
}

// Handler function to update an existing movie
func updateMovie(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var updatedMovie Movie

	// Bind the incoming JSON request to the Movie struct
	if err := c.Bind(&updatedMovie); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	updatedMovie.ID = id

	// Update the movie in the database
	if err := updateMovieInDB(updatedMovie); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedMovie)
}

// Handler function to delete a movie by ID
func deleteMovie(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	// Delete the movie from the database
	if err := deleteMovieFromDB(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
