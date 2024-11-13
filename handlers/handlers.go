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
	// Attempt to decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
		return
	}

	// Calculate points and respond with just the points
	points := service.CalculatePoints(&receipt)
	if points < 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Points{Points: points})
}
