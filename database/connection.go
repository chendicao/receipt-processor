package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv" // Import the godotenv package
	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect initializes the database connection
func Connect() error {
	// Load environment variables from .env file (if available)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Check if required environment variables are set
	if os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_NAME") == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	// Build the connection string
	connStr := "user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"

	// var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Set connection pool limits (optional but recommended)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Verify the connection is alive
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection is not alive: %v", err)
	}

	// Run database migrations
	runMigrations()

	return nil
}

// runMigrations creates tables if they don't exist
func runMigrations() {
	// Create receipts table if it doesn't exist
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS receipts (
			id UUID PRIMARY KEY,
			retailer TEXT NOT NULL,
			purchase_date DATE NOT NULL,
			purchase_time TIME NOT NULL,
			total NUMERIC NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create receipts table: %v", err)
	}

	// Create items table if it doesn't exist
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			receipt_id UUID REFERENCES receipts(id) ON DELETE CASCADE,
			short_description TEXT NOT NULL,
			price NUMERIC NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create items table: %v", err)
	}

	log.Println("Database tables created successfully!")
}

// Close closes the database connection
func Close() {
	DB.Close()
}
