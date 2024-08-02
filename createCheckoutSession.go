package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

// Declare the request structure for the line items and mode
type req struct {
	LineItems []LineItems `json:"line_items"`
	Mode      string      `json:"mode"` // Mode of the session (payment or subscription)
}

type LineItems struct {
	Title      string `json:"title"`
	Price      int64  `json:"price"`
	Quantity   int64  `json:"quantity"`
	CoverImage string `json:"coverImage"`
}

// Handler function to create a Stripe Checkout Session
func handleCreateCheckoutSession(c echo.Context) error {
	request := req{}
	// Bind the incoming JSON request to the req struct
	if err := c.Bind(&request); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
	}

	// Create parameters for the Stripe Checkout Session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(request.Mode),
		SuccessURL:         stripe.String("http://localhost:3000/success"),
		CancelURL:          stripe.String("http://localhost:3000/cancel"),
	}

	totalAmount := int64(0)
	// Add line items to the session parameters
	for _, item := range request.LineItems {
		productData := &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
			Name: stripe.String(item.Title),
		}

		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:    stripe.String("usd"),
				ProductData: productData,
				UnitAmount:  stripe.Int64(item.Price),
			},
			Quantity: stripe.Int64(item.Quantity),
		})
		totalAmount += item.Price * item.Quantity
	}

	// Create a new Stripe Checkout Session
	sess, err := session.New(params)
	if err != nil {
		log.Printf("Error creating Stripe session: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating Stripe session: " + err.Error()})
	}

	// Log the transaction
	err = logTransaction(sess.ID, int(totalAmount)/100, "created")
	if err != nil {
		log.Printf("Error logging transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error logging transaction: " + err.Error()})
	}

	// Return the session ID to the client
	resp := map[string]string{
		"id": sess.ID,
	}
	return c.JSON(http.StatusOK, resp)
}
