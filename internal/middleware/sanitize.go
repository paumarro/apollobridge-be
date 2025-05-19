package middleware

import (
	"html"
	"regexp"
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
				c.Request.Header[key][i] = sanitizeString(value)
			}
		}

		// Sanitize query parameters
		query := c.Request.URL.Query()
		for key, values := range query {
			for i, value := range values {
				query[key][i] = sanitizeString(value)
			}
		}
		c.Request.URL.RawQuery = query.Encode()

		// Sanitize path parameters
		for i, param := range c.Params {
			c.Params[i].Value = sanitizeString(param.Value)
		}

		// Sanitize body (if applicable)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			var req dto.ArtworkRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
				c.Abort()
				return
			}

			// Sanitize body fields
			req.Title = sanitizeString(req.Title)
			req.Artist = sanitizeString(req.Artist)
			req.Description = sanitizeString(req.Description)
			req.Image = sanitizeString(req.Image)

			// Save sanitized body to the context
			c.Set("sanitizedArtwork", req)
		}

		// Pass control to the next middleware/handler
		c.Next()
	}
}

// sanitizeString trims whitespace, escapes HTML, and removes HTML tags
func sanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)

	// Strip HTML tags
	stripped := stripHTMLTags(trimmed)

	// Escape HTML entities (e.g., < becomes &lt;, > becomes &gt;)
	escaped := html.EscapeString(stripped)
	return escaped
}

// stripHTMLTags removes HTML tags using a regular expression
func stripHTMLTags(input string) string {
	// Regex to match HTML tags
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(input, "")
}
