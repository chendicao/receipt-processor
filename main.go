package main

import (
	"log"
	"net/http"

	"github.com/chendicao/receipt-processor/database" // Import database for connection
	"github.com/chendicao/receipt-processor/handlers"
	"github.com/chendicao/receipt-processor/utils" // Import utils for middleware
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	database.Connect() // Establish the database connection

	// Initialize the router
	r := mux.NewRouter()

	// Apply middlewares
	r.Use(utils.RateLimiterMiddleware) // Rate-limiting middleware from utils
	r.Use(utils.CORSMiddleware())      // CORS middleware from middleware package (you could move this to utils if preferred)

	// Define routes
	r.HandleFunc("/receipts", handlers.ProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", handlers.GetReceiptPoints).Methods("GET")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
