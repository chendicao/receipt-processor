package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chendicao/receipt-processor/database" // Import the database package
	"github.com/chendicao/receipt-processor/models"   // Import the models package
	"github.com/chendicao/receipt-processor/service"  // Import the service package
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Initialize the validator
var validate = validator.New()

type ReceiptID struct {
	ID string `json:"id"`
}

type Points struct {
	Points int `json:"points"`
}

// ProcessReceipts handles receipt processing and stores it in the database
func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// Validate receipt structure
	if err := validate.Struct(receipt); err != nil {
		http.Error(w, "Invalid receipt structure: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a UUID for the receipt ID
	receiptID := uuid.New().String()

	// Start a transaction to insert the receipt and its items
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to begin transaction:", err)
		return
	}
	defer tx.Rollback() // Ensure rollback if anything fails

	// Insert receipt into the database
	_, err = tx.Exec(`
		INSERT INTO receipts (id, retailer, purchase_date, purchase_time, total)
		VALUES ($1, $2, $3, $4, $5)
	`, receiptID, receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total)
	if err != nil {
		http.Error(w, "Failed to insert receipt", http.StatusInternalServerError)
		log.Println("Failed to insert receipt:", err)
		return
	}

	// Insert items into the database
	for _, item := range receipt.Items {
		_, err := tx.Exec(`
			INSERT INTO items (receipt_id, short_description, price)
			VALUES ($1, $2, $3)
		`, receiptID, item.ShortDescription, item.Price)
		if err != nil {
			http.Error(w, "Failed to insert item", http.StatusInternalServerError)
			log.Println("Failed to insert item:", err)
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		log.Println("Failed to commit transaction:", err)
		return
	}

	// Respond with only the generated ID
	response := ReceiptID{ID: receiptID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetReceiptPoints handles fetching receipt points based on the receipt ID
func GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the receipt from the database
	var receipt models.Receipt
	err := database.DB.QueryRow(`
		SELECT retailer, purchase_date, purchase_time, total 
		FROM receipts WHERE id = $1
	`, id).Scan(&receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No receipt found for that id", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Println("Failed to retrieve receipt:", err)
		return
	}

	// Retrieve items for the receipt
	rows, err := database.DB.Query(`
		SELECT short_description, price FROM items WHERE receipt_id = $1
	`, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to retrieve items:", err)
		return
	}
	defer rows.Close()

	// Add items to the receipt
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ShortDescription, &item.Price); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to scan item:", err)
			return
		}
		receipt.Items = append(receipt.Items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Failed to process items:", err)
		return
	}

	// Calculate points and respond with just the points
	points := service.CalculatePoints(&receipt)
	if points < 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return points as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Points{Points: points})
}
