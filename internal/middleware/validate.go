package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/paumarro/apollo-be/internal/dto"
)

// Validator instance
var validate = validator.New()

// Validate middleware validates the sanitized input
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate query parameters
		if err := validateQueryParams(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Validate path parameters
		if err := validatePathParams(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Validate body (if applicable)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if err := validateRequestBody(c); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err})
				c.Abort()
				return
			}
		}

		// Pass control to the next middleware/handler
		c.Next()
	}
}

// validateQueryParams validates query parameters
func validateQueryParams(c *gin.Context) error {
	query := c.Request.URL.Query()
	for key, values := range query {
		if len(key) > 50 {
			return fmt.Errorf("query parameter key '%s' is too long", key)
		}
		for _, value := range values {
			decodedValue, err := url.QueryUnescape(value)
			if err != nil {
				return fmt.Errorf("failed to decode query parameter '%s'", key)
			}

			if len(decodedValue) > 255 {
				return fmt.Errorf("query parameter value for key '%s' is too long", key)
			}
			if !isSafeString(decodedValue) {
				return fmt.Errorf("query parameter value for key '%s' contains invalid characters", key)
			}
		}
	}
	return nil
}

// validatePathParams validates path parameters
func validatePathParams(c *gin.Context) error {
	for _, param := range c.Params {
		if len(param.Value) > 255 {
			return fmt.Errorf("path parameter '%s' is too long", param.Key)
		}
		if !isSafeString(param.Value) {
			return fmt.Errorf("path parameter '%s' contains invalid characters", param.Key)
		}
	}
	return nil
}

// validateRequestBody validates the request body
func validateRequestBody(c *gin.Context) map[string]string {
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		return map[string]string{"error": "Input not sanitized"}
	}

	req := sanitizedArtwork.(dto.ArtworkRequest)

	// Validate the sanitized input
	if err := validate.Struct(req); err != nil {
		return formatValidationErrors(err)
	}

	return nil
}

// isSafeString validates a string using a whitelist approach
func isSafeString(input string) bool {
	// Allow alphanumeric characters, spaces, and a few safe symbols
	allowed := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,@!\/]*$`)
	return allowed.MatchString(input)
}

// formatValidationErrors formats validation errors into a user-friendly map
func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		fieldName := fieldErr.Field()
		tag := fieldErr.Tag()

		// Customize error messages based on validation tags
		var message string
		switch tag {
		case "required":
			message = "This field is required."
		case "max":
			message = "This field exceeds the maximum allowed length."
		case "min":
			message = "This field is below the minimum required length."
		case "email":
			message = "Invalid email address format."
		default:
			message = "Invalid value."
		}

		errors[fieldName] = message
	}
	return errors
}
