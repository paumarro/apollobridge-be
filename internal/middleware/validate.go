package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/paumarro/apollo-be/internal/dto"
)

// Validator instance
var validate = validator.New()

// Validate middleware validates the sanitized input
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodDelete {
			c.Next() // Pass control to the next middleware/handler
			return
		}
		// Retrieve the sanitized input from the context
		sanitizedArtwork, exists := c.Get("sanitizedArtwork")
		if !exists {
			c.JSON(400, gin.H{"error": "Input not sanitized"})
			c.Abort()
			return
		}
		req := sanitizedArtwork.(dto.ArtworkRequest)

		// Validate the sanitized input
		if err := validate.Struct(req); err != nil {
			c.JSON(400, gin.H{"error": "Validation failed", "details": formatValidationErrors(err)})
			c.Abort()
			return
		}

		// Pass control to the next middleware/handler
		c.Next()
	}
}

// Helper function to format validation errors
func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		errors[fieldErr.Field()] = fieldErr.Tag() // Example: "Title": "required"
	}
	return errors
}
