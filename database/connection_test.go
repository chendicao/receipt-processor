package database

import (
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// Set environment variables for testing (you may use a testing .env file or mock them)
	// os.Setenv("DB_USER", "your_user")
	// os.Setenv("DB_PASSWORD", "your_password")
	// os.Setenv("DB_NAME", "test_db")
	err := godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}
	// Connect to the database
	err = Connect()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test that the database connection is live
	err = DB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping the database: %v", err)
	}

	// Assert that DB is not nil
	assert.NotNil(t, DB)

	// Clean up after the test by closing the DB connection
	defer Close()
}

func TestRunMigrations(t *testing.T) {
	// Set environment variables for testing
	// os.Setenv("DB_USER", "your_user")
	// os.Setenv("DB_PASSWORD", "your_password")
	// os.Setenv("DB_NAME", "test_db")
	err := godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}
	// Connect to the database
	err = Connect()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Run the migrations
	runMigrations()

	// Check if the tables were created (e.g., receipts and items tables)
	var tableCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'receipts'").Scan(&tableCount)
	if err != nil {
		t.Fatalf("Failed to check the receipts table: %v", err)
	}
	assert.Equal(t, 1, tableCount)

	// Clean up after the test by closing the DB connection
	defer Close()
}
