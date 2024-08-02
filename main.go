package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stripe/stripe-go/v78"
)

func main() {

	// Initialize Stripe with your secret key
	stripe.Key = "sk_test_51PM3XgJ8mel0IwGSfryYCLjvac4oQfJgYxApA3rYDa94gmLraOqNUrYgT1iA4rbb86Ygh7zNHQPrm8acF0jgZSKf00JWQC5EEY"

	// Initialize the database
	initDB()

	// Ensure the database is closed when the application stops
	defer db.Close()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Middleware for handling Cross-Origin Resource Sharing (CORS)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	// Serve static files from "public/images" directory
	e.Static("/images", "public/images")

	// Register movie routes
	e.GET("/movies", getMovies)
	e.GET("/movies/:id", getMovie)
	e.POST("/movies", createMovie)
	e.PUT("/movies/:id", updateMovie)
	e.DELETE("/movies/:id", deleteMovie)

	// Register cart routes
	e.GET("/cart", getCart)
	e.POST("/cart", addToCart)
	e.DELETE("/cart/:id", removeFromCart)

	// Register auth routes
	e.POST("/signup", signup)
	e.POST("/login", login)

	// Protected route
    e.GET("/protected-endpoint", protectedEndpoint, jwtMiddleware)

	// Register endpoint handler
	e.POST("/create-checkout-session", handleCreateCheckoutSession, jwtMiddleware)

	// Register transaction history route
	e.GET("/transactions", getAllTransactions, jwtMiddleware)

	// Log server start message
	log.Println("Server listening on :8080")

	// Start the server on port 8080 and log any fatal errors
	log.Fatal(e.Start(":8080"))
}
