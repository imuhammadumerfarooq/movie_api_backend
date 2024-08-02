package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Transaction struct definition
type Transaction struct {
	ID        int    `json:"id"`
	SessionID string `json:"session_id"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// Function to log a new transaction to the database
func logTransaction(sessionID string, amount int, status string) error {
	_, err := db.Exec("INSERT INTO transactions (session_id, amount, status) VALUES (?, ?, ?)", sessionID, amount, status)
	return err
}

// Handler function to get all transactions
func getAllTransactions(c echo.Context) error {
	rows, err := db.Query("SELECT id, session_id, amount, status, created_at FROM transactions")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching transactions"})
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.SessionID, &transaction.Amount, &transaction.Status, &transaction.CreatedAt); err != nil {
			log.Printf("Error scanning transaction: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning transaction"})
		}
		log.Printf("Transaction retrieved: %+v", transaction)
		transactions = append(transactions, transaction)
	}

	return c.JSON(http.StatusOK, transactions)
}
