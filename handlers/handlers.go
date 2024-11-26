package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chendicao/receipt-processor/models"
	"github.com/chendicao/receipt-processor/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	receiptStore       = make(map[string]models.Receipt) // Stores receipts by receipt ID
	receiptPointsStore = make(map[string]Points)         // Stores points for each receipt
	userStore          = make(map[string]UserData)       // Stores user data (points, receipt count)
	validate           = validator.New()                 // Used for validating the structure of the receipt
)

type ReceiptID struct {
	ID string `json:"id"`
}

type Points struct {
	Points int `json:"points"`
}

type UserData struct {
	TotalPoints  int `json:"total_points"`  // Total points for the user
	ReceiptCount int `json:"receipt_count"` // Number of receipts uploaded by the user
}

// ProcessReceipts handles the receipt upload and calculation of points.
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

	// Calculate base points for the receipt using the CalculatePoints method
	basePoints := service.CalculatePoints(&receipt)

	// Retrieve user data and calculate bonus points
	userData, exists := userStore[receipt.UserID]
	if !exists {
		userData = UserData{
			TotalPoints:  0,
			ReceiptCount: 0,
		}
	}

	// Calculate bonus points based on the receipt count (order)
	bonusPoints := 0
	if userData.ReceiptCount > 0 {
		bonusPoints = userData.ReceiptCount * 10 // Bonus points start from the second receipt
	}

	// Total points for the current receipt
	totalPoints := basePoints + bonusPoints

	// Save the receipt in the receipt store
	receiptStore[receiptID] = receipt

	// Save the points for this receipt in a separate map (not updating user points yet)
	receiptPoints := Points{Points: totalPoints}
	receiptPointsStore[receiptID] = receiptPoints // Store points for each receipt separately

	// Update the user's total points and increment receipt count
	userData.TotalPoints += totalPoints
	userData.ReceiptCount++

	// Save updated user data
	userStore[receipt.UserID] = userData

	// Respond with the generated receipt ID
	response := ReceiptID{ID: receiptID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetReceiptPoints handles retrieving the points for a specific receipt.
func GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the points for the specific receipt from the store
	receiptPoints, ok := receiptPointsStore[id]
	if !ok {
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
		return
	}

	// Respond with the points for the specific receipt
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receiptPoints)
}

// GetUserPoints handles retrieving the total points for a specific user.
func GetUserPoints(w http.ResponseWriter, r *http.Request) {
	// Get the user_id from the route
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Check if user data exists
	userData, exists := userStore[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user's total points
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Points{Points: userData.TotalPoints})
}
