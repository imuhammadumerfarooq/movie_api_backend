package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Global user variable to store user data
var user User

// User struct definition
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handler function to sign up a new user
func signup(c echo.Context) error {

	// Bind the incoming JSON request to the User struct
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate input fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username, email, and password are required"})
	}

	// Check if email already exists in the database
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email)
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error checking email:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	// Return a bad request status if the email already exists
	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email already exists"})
	}

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	// Insert the new user into the database
	insertUserSQL := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err = db.Exec(insertUserSQL, user.Username, user.Email, string(hashedPassword))
	if err != nil {
		log.Println("Error inserting user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	// Return a created status with a success message
	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
