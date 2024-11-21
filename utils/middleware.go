package utils

import (
	"net/http"
	"time"

	"github.com/rs/cors" // Add the missing CORS import
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory" // Use the memory store instead of redis
)

var limiterInstance *limiter.Limiter

// Initialize in-memory store and rate limiter
func init() {
	// Use the memory store for rate limiting
	store := memory.NewStore() // In-memory store

	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  10, // Max 10 requests per second
	}

	limiterInstance = limiter.New(store, rate)
}

// RateLimiterMiddleware to limit requests
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the context and check if the rate limit is reached
		ctx, err := limiterInstance.Get(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if ctx.Reached {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// If the rate limit isn't reached, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware to handle CORS configuration
func CORSMiddleware() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Use specific origins in production for security
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true, // Ensure this is needed for your use case
	}).Handler
}
