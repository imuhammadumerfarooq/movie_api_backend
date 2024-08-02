package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CartItem represents an item in the cart
type CartItem struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	Genre      string `json:"genre"`
	Rating     string `json:"rating"`
	CoverImage string `json:"coverImage"`
	Quantity   int    `json:"quantity"`
}

// Global cart variable to store cart items
var cart []CartItem

// Handler function to get all items in the cart
func getCart(c echo.Context) error {
	return c.JSON(http.StatusOK, cart)
}

// Handler function to add an item to the cart
func addToCart(c echo.Context) error {
	var item CartItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// Check if the item is already in the cart
	for i, movie := range cart {
		if movie.ID == item.ID {
			cart[i].Quantity += item.Quantity
			return c.JSON(http.StatusOK, cart)
		}
	}

	// Add the new item to the cart if it is not already present
	cart = append(cart, item)

	return c.JSON(http.StatusOK, cart)
}

// Handler function to remove an item from the cart by ID
func removeFromCart(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	// Iterate over the cart items to find the one with the matching ID
	for i, movie := range cart {
		if movie.ID == id {
			cart = append(cart[:i], cart[i+1:]...)
			return c.JSON(http.StatusOK, cart)
		}
	}

	return c.JSON(http.StatusNotFound, "Item not found in cart")
}
