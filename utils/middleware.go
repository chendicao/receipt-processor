package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/cors" // Add the missing CORS import
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

var limiterInstance *limiter.Limiter
var redisClient *redis.Client

// Initialize Redis client and rate limiter
func init() {
	// Create Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Default DB
	})
	// Ensure Redis is reachable
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}

	// Use Redis as the store for rate limiting
	store := redis.NewStore(redisClient)
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
