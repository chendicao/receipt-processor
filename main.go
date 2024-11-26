package main

import (
	"log"
	"net/http"

	"github.com/chendicao/receipt-processor/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the router
	router := mux.NewRouter()

	// Define your endpoints and handlers
	router.HandleFunc("/receipts/process", handlers.ProcessReceipts).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handlers.GetReceiptPoints).Methods("GET")
	router.HandleFunc("/users/{user_id}/points", handlers.GetUserPoints).Methods("GET") // New route for user points

	// Start the server
	log.Println("Starting the server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
