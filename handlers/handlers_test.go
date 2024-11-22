package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chendicao/receipt-processor/database"
	"github.com/chendicao/receipt-processor/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Mock database initialization
func setupMockDB() (sqlmock.Sqlmock, *sql.DB, func()) {
	// Open a mock database
	db, mock, _ := sqlmock.New()

	// This function will clean up after the test
	return mock, db, func() { db.Close() }
}

func TestProcessReceipts(t *testing.T) {
	// Setup mock database
	mockDB, db, cleanup := setupMockDB()
	defer cleanup()

	// Assign the mock database to the global DB variable
	database.DB = db

	// Mock expected database behavior
	mockDB.ExpectExec("INSERT INTO receipts").WithArgs(
		"Test Retailer", "2024-11-22", "12:34", 100.00,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new request to pass into the handler
	receipt := models.Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2024-11-22",
		PurchaseTime: "12:34",
		Total:        100.00,
	}

	// Convert the receipt to JSON
	receiptJSON, _ := json.Marshal(receipt)

	req, err := http.NewRequest("POST", "/receipts", bytes.NewBuffer(receiptJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create the router and apply the handler
	r := mux.NewRouter()
	r.HandleFunc("/receipts", ProcessReceipts).Methods("POST")

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	assert.Contains(t, response["id"], "uuid") // Check if the response contains an ID (UUID)
}

func TestGetReceiptPoints(t *testing.T) {
	// Setup mock database
	mockDB, db, cleanup := setupMockDB()
	defer cleanup()

	// Assign the mock database to the global DB variable
	database.DB = db

	// Mock expected database behavior
	mockDB.ExpectQuery("SELECT retailer, purchase_date, purchase_time, total").
		WithArgs("test-id").
		WillReturnRows(sqlmock.NewRows([]string{"retailer", "purchase_date", "purchase_time", "total"}).
			AddRow("Test Retailer", "2024-11-22", "12:34", 100.00))

	// Create a new request to pass into the handler
	req, err := http.NewRequest("GET", "/receipts/test-id/points", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create the router and apply the handler
	r := mux.NewRouter()
	r.HandleFunc("/receipts/{id}/points", GetReceiptPoints).Methods("GET")

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]int
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	assert.Equal(t, 10, response["points"]) // Assuming the points calculation returns 10
}
