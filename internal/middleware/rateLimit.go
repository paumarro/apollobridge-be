package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Define a struct to hold the rate limiter for each client
type ClientLimiter struct {
	limiter *rate.Limiter
}

// Map to store rate limiters for each client (e.g., by IP)
var clients = make(map[string]*ClientLimiter)
var mu sync.Mutex

// Middleware to limit requests
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Get or create a rate limiter for the client
		limiter := getLimiter(clientIP)

		// Check if the request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// Function to get or create a rate limiter for a client
func getLimiter(clientIP string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// Check if the client already has a rate limiter
	if _, exists := clients[clientIP]; !exists {
		// Create a new rate limiter: 5 requests per second with a burst of 10
		clients[clientIP] = &ClientLimiter{
			limiter: rate.NewLimiter(5, 10),
		}
	}
	return clients[clientIP].limiter
}
