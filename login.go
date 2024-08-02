package main

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Secret key used to sign JWT tokens
var jwtKey = []byte("1SmqmurUmjhcLE2o0VXJWH7sfL3Ibh3S")

// Claims struct for JWT
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Define input struct for email and password
var input struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handler function to sign in a user
func login(c echo.Context) error {

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Retrieve the user from the database by email
	var user User
	row := db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", input.Email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error fetching user:", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Compare the hashed password with the password provided
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Generate JWT token with a 24-hour expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}


