package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Global database connection pool
var db *sql.DB

// Function to initialize the database
func initDB() {
	var err error

	// Open a connection to the SQLite database
	db, err = sql.Open("sqlite3", "./movies.db")
	if err != nil {
		// Log a fatal error if the database cannot be opened
		log.Fatalf("Error opening database: %v", err)
	}

	// SQL statement to create the Movie table if it does not exist
	createMovieTable := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY,
		title TEXT,
		year INTEGER,
		genre TEXT,
		rating TEXT,
		coverImage TEXT,
		price INTEGER
	);
	`

	// Execute the SQL statement to create the Movie table
	_, err = db.Exec(createMovieTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// SQL statement to create the Users table if it does not exist
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`

	// Execute the SQL statement to create the Users table
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	// SQL statement to create the Transactions table if it does not exist
	createTransactionsTable := `
    CREATE TABLE IF NOT EXISTS transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        session_id TEXT,
        amount INTEGER,
        status TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

`

	// Execute the SQL statement to create the Transactions table
	_, err = db.Exec(createTransactionsTable)
	if err != nil {
		log.Fatalf("Error creating transactions table: %v", err)
	}
}
