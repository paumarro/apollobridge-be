package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/dto"
	"golang.org/x/text/unicode/norm"
)

// Sanitize strictly parses and normalizes ArtworkRequest bodies for POST/PUT.
// - Enforces Content-Type: application/json
// - Limits body size (1 MB)
// - Rejects unknown JSON fields (DisallowUnknownFields)
// - Normalizes strings (TrimSpace + Unicode NFC) without HTML escaping or tag stripping
// - Stores the normalized request in context as "sanitizedArtwork"
func Sanitize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only handle JSON bodies for POST/PUT. Do not mutate headers/query/path.
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut {
			c.Next()
			return
		}

		// Enforce Content-Type
		ct := c.GetHeader("Content-Type")
		if ct == "" || !strings.HasPrefix(strings.ToLower(ct), "application/json") {
			errorResponse(c, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json", nil)
			return
		}

		// Limit body size to 1MB (adjust as needed)
		const maxBodyBytes = 1 << 20 // 1MB
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

		// Strict decode with DisallowUnknownFields
		var req dto.ArtworkRequest
		dec := json.NewDecoder(c.Request.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			// Map common errors
			msg, code, details, status := mapJSONDecodeError(err)
			errorResponse(c, status, code, msg, details)
			return
		}

		// Normalize fields (no HTML escaping or tag stripping)
		req.Title = normalizeString(req.Title)
		req.Artist = normalizeString(req.Artist)
		req.Description = normalizeString(req.Description)
		req.Image = normalizeString(req.Image)

		// Re-encode normalized JSON back into the request body for any downstream binders (optional)
		buf, _ := json.Marshal(req)
		c.Request.Body = io.NopCloser(bytes.NewReader(buf))

		// Store normalized request in context for controllers/validators
		c.Set("sanitizedArtwork", req)

		c.Next()
	}
}

func normalizeString(s string) string {
	// Unicode NFC normalization plus trim. Avoid HTML escaping here; do output encoding at render time.
	s = norm.NFC.String(s)
	s = strings.TrimSpace(s)
	return s
}

func mapJSONDecodeError(err error) (message, code string, details map[string]string, status int) {
	e := err.Error()
	// Payload too large comes from MaxBytesReader
	if strings.Contains(e, "http: request body too large") {
		return "Payload too large", "payload_too_large", nil, http.StatusRequestEntityTooLarge
	}
	// Unknown field: json: unknown field "..."
	if strings.Contains(e, "unknown field") {
		field := extractUnknownFieldName(e)
		details = map[string]string{}
		if field != "" {
			details["unknown_fields"] = field
		}
		return "Unknown field in JSON payload", "unknown_field", details, http.StatusBadRequest
	}
	// Generic invalid JSON
	return "Invalid JSON payload", "invalid_json", nil, http.StatusBadRequest
}

func extractUnknownFieldName(errStr string) string {
	// Typical form: "json: unknown field \"date\""
	start := strings.Index(errStr, "\"")
	end := strings.LastIndex(errStr, "\"")
	if start >= 0 && end > start {
		return errStr[start+1 : end]
	}
	return ""
}
