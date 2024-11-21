package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// Check if environment variables are set
	if os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_NAME") == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	connStr := "user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Set connection pool limits (optional but recommended)
	DB.SetMaxOpenConns(10) // Set maximum open connections
	DB.SetMaxIdleConns(5)  // Set maximum idle connections

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection is not alive: %v", err)
	}
}

func Close() {
	DB.Close() // No need to check for error here
}
