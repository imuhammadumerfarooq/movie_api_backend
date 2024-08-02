package main

import (
	"fmt"
)

// Function to read all movies from the database
func readMovies(limit, offset int) ([]Movie, error) {

	query := ("SELECT id, title, year, genre, rating, coverImage, price FROM movies LIMIT ? OFFSET ?")
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre, &movie.Rating, &movie.Image, &movie.Price); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func readMovieById(id string) (*Movie, error) {
	query := fmt.Sprintf("SELECT id, title, year, genre, rating, coverImage, price FROM movies where id = %v", id)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movie Movie
	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre, &movie.Rating, &movie.Image, &movie.Price); err != nil {
			return nil, err
		}
	}

	return &movie, nil
}

// Function to write a new movie to the database
func writeMovie(movie Movie) error {
	_, err := db.Exec("INSERT INTO movies (title, year, genre, rating, coverImage, price) VALUES (?, ?, ?, ?, ?, ?)", movie.Title, movie.Year, movie.Genre, movie.Rating, &movie.Image, movie.Price)
	return err
}

// Function to update an existing movie in the database
func updateMovieInDB(movie Movie) error {
	_, err := db.Exec("UPDATE movies SET title = ?, year = ?, genre = ?, rating = ?, coverImage = ?, price = ? WHERE id = ?", movie.Title, movie.Year, movie.Genre, movie.Rating, &movie.Image, movie.Price, movie.ID)
	return err
}

// Function to delete a movie from the database
func deleteMovieFromDB(id int) error {
	_, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	return err
}
