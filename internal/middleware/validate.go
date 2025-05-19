package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/paumarro/apollo-be/internal/dto"
)

var validate = validator.New()

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateQueryParams(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if err := validatePathParams(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if err := validateRequestBody(c); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func validateQueryParams(c *gin.Context) error {
	query := c.Request.URL.Query()
	for key, values := range query {
		for _, value := range values {
			switch key {
			case "state":
				if !strings.HasPrefix(value, "/") {
					return fmt.Errorf("invalid 'state' parameter: must start with '/'")
				}
			case "session_state":
				if !govalidator.IsUUID(value) {
					return fmt.Errorf("invalid 'session_state': must be a UUID")
				}
			case "iss":
				decoded, err := url.QueryUnescape(value)
				if err != nil || !govalidator.IsURL(decoded) || !strings.HasPrefix(decoded, "https://") {
					return fmt.Errorf("invalid 'iss' parameter: must be a valid HTTPS URL")
				}
			case "code":
				if len(value) < 20 || !regexp.MustCompile(`^[a-zA-Z0-9\.\-]+$`).MatchString(value) {
					return fmt.Errorf("invalid 'code' parameter: must be alphanumeric with dots/hyphens")
				}
			default:
				if !govalidator.IsPrintableASCII(value) {
					return fmt.Errorf("query parameter '%s' contains unsafe characters", key)
				}
			}
		}
	}
	return nil
}

func validatePathParams(c *gin.Context) error {
	for _, param := range c.Params {
		if len(param.Value) > 255 {
			return fmt.Errorf("path parameter '%s' is too long", param.Key)
		}
		if !govalidator.IsPrintableASCII(param.Value) {
			return fmt.Errorf("path parameter '%s' contains invalid characters", param.Key)
		}
	}
	return nil
}

func validateRequestBody(c *gin.Context) map[string]string {
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		return map[string]string{"error": "Input not sanitized"}
	}

	req := sanitizedArtwork.(dto.ArtworkRequest)

	if err := validate.Struct(req); err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		field := fieldErr.Field()
		tag := fieldErr.Tag()
		var message string

		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required.", field)
		case "max":
			message = fmt.Sprintf("%s exceeds the maximum allowed length.", field)
		case "min":
			message = fmt.Sprintf("%s is below the minimum required length.", field)
		case "email":
			message = fmt.Sprintf("%s is not a valid email address.", field)
		default:
			message = fmt.Sprintf("%s has an invalid value.", field)
		}

		errors[field] = message
	}
	return errors
}
