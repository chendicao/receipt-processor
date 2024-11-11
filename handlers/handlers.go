package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chendicao/receipt-processor/models"
	"github.com/chendicao/receipt-processor/service"
	"github.com/google/uuid" // Import the UUID package
	"github.com/gorilla/mux"
)

var receiptStore = make(map[string]models.Receipt)

type ReceiptID struct {
	ID string `json:"id"`
}

type Points struct {
	Points int `json:"points"`
}

func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	_ = json.NewDecoder(r.Body).Decode(&receipt)

	// Generate a UUID for the receipt ID
	receiptID := uuid.New().String()

	// Store the receipt in memory with the generated UUID as the key
	receiptStore[receiptID] = receipt

	// Respond with only the generated ID
	response := ReceiptID{ID: receiptID}
	json.NewEncoder(w).Encode(response)
}

func GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the receipt from the store
	receipt, ok := receiptStore[id]
	if !ok {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Calculate points and respond with just the points
	points := service.CalculatePoints(&receipt)
	json.NewEncoder(w).Encode(Points{Points: points})
}
