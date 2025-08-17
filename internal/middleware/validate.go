package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/paumarro/apollo-be/internal/dto"
)

var validate = validator.New()

// Validate performs request-level validation without mutating inputs.
// - Validates Accept header (if provided) for JSON.
// - Validates known query parameters for auth/OIDC-like flows.
// - Validates the ArtworkRequest body (POST/PUT) after Sanitize() parsed and normalized it.
// - Enforces HTTPS for image URLs and rejects control characters where disallowed.
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Accept header validation (lenient if header missing)
		if err := validateAcceptHeader(c); err != nil {
			errorResponse(c, http.StatusNotAcceptable, "not_acceptable", "Requested representation not acceptable", nil)
			c.Abort()
			return
		}

		if err := validateQueryParams(c); err != nil {
			// Map query param violations to validation_error
			details := map[string]string{"query": err.Error()}
			errorResponse(c, http.StatusBadRequest, "validation_error", "Request validation failed", details)
			c.Abort()
			return
		}

		// Only validate body for POST/PUT if Sanitize() already set it.
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			if details := validateRequestBody(c); details != nil {
				errorResponse(c, http.StatusBadRequest, "validation_error", "Request validation failed", details)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func validateAcceptHeader(c *gin.Context) error {
	// Be lenient: if absent, allow. If present, must allow JSON.
	accept := c.GetHeader("Accept")
	if strings.TrimSpace(accept) == "" {
		return nil
	}
	// Accept may be a list. Allow */*, application/json, application/*+json
	parts := strings.Split(accept, ",")
	for _, p := range parts {
		mt := strings.ToLower(strings.TrimSpace(strings.Split(p, ";")[0]))
		if mt == "*/*" || mt == "application/json" || (strings.HasPrefix(mt, "application/") && strings.HasSuffix(mt, "+json")) {
			return nil
		}
	}
	return fmt.Errorf("not acceptable")
}

func validateQueryParams(c *gin.Context) error {
	q := c.Request.URL.Query()

	// state: must begin with "/", reasonable length
	if v := q.Get("state"); v != "" {
		if !strings.HasPrefix(v, "/") || len(v) > 2048 {
			return fmt.Errorf("invalid 'state' parameter")
		}
	}

	// session_state: UUID
	if v := q.Get("session_state"); v != "" {
		if !govalidator.IsUUID(v) {
			return fmt.Errorf("invalid 'session_state'")
		}
	}

	// iss: must be a valid HTTPS URL
	if v := q.Get("iss"); v != "" {
		u, err := url.Parse(v)
		if err != nil || u.Scheme != "https" || u.Host == "" {
			return fmt.Errorf("invalid 'iss' parameter")
		}
	}

	// code: base64url-safe-like allowlist, length-bound
	if v := q.Get("code"); v != "" {
		if len(v) < 20 || len(v) > 512 {
			return fmt.Errorf("invalid 'code' parameter")
		}
		if !regexp.MustCompile(`^[A-Za-z0-9._~\-]+$`).MatchString(v) {
			return fmt.Errorf("invalid 'code' parameter")
		}
	}

	return nil
}

func validateRequestBody(c *gin.Context) map[string]string {
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		return map[string]string{"body": "input not sanitized"}
	}

	req, ok := sanitizedArtwork.(dto.ArtworkRequest)
	if !ok {
		return map[string]string{"body": "invalid sanitized input"}
	}

	// Struct tag validation
	if err := validate.Struct(req); err != nil {
		return formatValidationErrors(err)
	}

	// Enforce HTTPS for image
	if msg := requireHTTPS(req.Image); msg != "" {
		return map[string]string{"image": msg}
	}

	// Control character checks
	if containsControlChars(req.Title, false) {
		return map[string]string{"title": "title contains invalid control characters"}
	}
	if containsControlChars(req.Artist, false) {
		return map[string]string{"artist": "artist contains invalid control characters"}
	}
	// Description: allow \n, \r, \t but no other control characters
	if containsControlChars(req.Description, true) {
		return map[string]string{"description": "description contains invalid control characters"}
	}
	if containsControlChars(req.Image, false) {
		return map[string]string{"image": "image contains invalid control characters"}
	}

	return nil
}

func requireHTTPS(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return "image must be a valid https URL"
	}
	return ""
}

func containsControlChars(s string, allowNewlines bool) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			switch r {
			case '\n', '\r', '\t':
				if allowNewlines {
					continue
				}
			}
			return true
		}
	}
	return false
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if fe, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range fe {
			field := strings.ToLower(fieldErr.Field())
			switch fieldErr.Tag() {
			case "required":
				errors[field] = field + " is required."
			case "max":
				errors[field] = field + " exceeds the maximum allowed length."
			case "min":
				errors[field] = field + " is below the minimum required length."
			case "url":
				errors[field] = field + " must be a valid URL."
			default:
				errors[field] = field + " has an invalid value."
			}
		}
		return errors
	}

	// Fallback
	errors["error"] = "invalid request"
	return errors
}

// errorResponse builds a consistent error payload.
func errorResponse(c *gin.Context, status int, code, message string, details map[string]string) {
	payload := gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	}
	if len(details) > 0 {
		payload["error"].(gin.H)["details"] = details
	}
	c.AbortWithStatusJSON(status, payload)
}
