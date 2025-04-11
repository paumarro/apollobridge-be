package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/dto"
)

// Sanitize middleware sanitizes headers, query parameters, path parameters, and body data
func Sanitize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize headers
		for key, values := range c.Request.Header {
			for i, value := range values {
				c.Request.Header[key][i] = strings.TrimSpace(value) // Trim whitespace
			}
		}

		// Sanitize query parameters
		query := c.Request.URL.Query()
		for key, values := range query {
			for _, value := range values {
				query.Set(key, strings.TrimSpace(value)) // Trim whitespace
			}
		}
		c.Request.URL.RawQuery = query.Encode() // Update the query string

		// Sanitize path parameters
		for i, param := range c.Params {
			c.Params[i].Value = strings.TrimSpace(param.Value) // Trim whitespace
		}

		// Sanitize body (if applicable)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			var req dto.ArtworkRequest
			if err := c.ShouldBindJSON(&req); err == nil {
				// Trim whitespace in body fields
				req.Title = strings.TrimSpace(req.Title)
				req.Artist = strings.TrimSpace(req.Artist)
				req.Description = strings.TrimSpace(req.Description)
				req.Image = strings.TrimSpace(req.Image)

				// Save sanitized body to the context
				c.Set("sanitizedArtwork", req)
			}
		}

		// Pass control to the next middleware/handler
		c.Next()
	}
}
