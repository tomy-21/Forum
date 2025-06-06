package config

import (
	"database/sql"
	"fmt" // Import fmt for error wrapping
	"log"
	"os"

	"github.com/joho/godotenv" // go get github.com/joho/godotenv
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
}

// InitDB initializes and returns the database connection.
// This version directly opens the connection based on DSN from env
// and assumes you will handle table creation separately or that
// your database.Init/CreateTables functions will be called with the *sql.DB instance.
func InitDB() (*sql.DB, error) {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		return nil, fmt.Errorf("DB_DSN not set in environment variables")
	}

	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to MySQL database!")
	// Table creation should be handled explicitly after this,
	// for example by calling a refactored database.CreateTables(db)
	return db, nil
}
