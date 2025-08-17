package unit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/middleware"
	"github.com/stretchr/testify/require"
)

func setupSanitizeRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.Sanitize())
	// Dummy endpoint returns sanitizedArtwork to inspect normalization and body replacement
	r.POST("/gallery/artworks", func(c *gin.Context) {
		if v, ok := c.Get("sanitizedArtwork"); ok {
			c.JSON(http.StatusOK, gin.H{"artwork": v})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing sanitizedArtwork"})
	})
	return r
}

func TestSanitize_UnsupportedMediaType(t *testing.T) {
	r := setupSanitizeRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "text/plain")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnsupportedMediaType, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "unsupported_media_type", errObj["code"])
	require.Equal(t, "Content-Type must be application/json", errObj["message"])
}

func TestSanitize_PayloadTooLarge(t *testing.T) {
	r := setupSanitizeRouter()

	// Build a body > 1MB
	largeDesc := strings.Repeat("x", (1<<20)+100) // 1MB + 100
	payload := map[string]any{
		"title":       "Valid Title",
		"artist":      "Valid Artist",
		"description": largeDesc,
		"image":       "https://example.com/i.png",
	}
	b, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusRequestEntityTooLarge, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "payload_too_large", errObj["code"])
	require.Equal(t, "Payload too large", errObj["message"])
}

func TestSanitize_InvalidJSON(t *testing.T) {
	r := setupSanitizeRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "invalid_json", errObj["code"])
	require.Equal(t, "Invalid JSON payload", errObj["message"])
}

func TestSanitize_UnknownFieldRejected(t *testing.T) {
	r := setupSanitizeRouter()

	payload := map[string]any{
		"title":       "Valid Title",
		"artist":      "Valid Artist",
		"description": "desc",
		"image":       "https://example.com/i.png",
		"date":        "2020-01-01T00:00:00Z",
	}
	b, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "unknown_field", errObj["code"])
	require.Equal(t, "Unknown field in JSON payload", errObj["message"])
	// details.unknown_fields should be "date"
	if details, ok := errObj["details"].(map[string]any); ok {
		require.Equal(t, "date", details["unknown_fields"])
	} else {
		t.Fatalf("expected details in error response")
	}
}

func TestSanitize_NormalizesAndStores(t *testing.T) {
	r := setupSanitizeRouter()

	// Artist uses A + combining acute; NFC normalization should yield "Á"
	payload := map[string]any{
		"title":       "  Valid Title  ",
		"artist":      "A\u0301rtist", // Ártist expected after normalization
		"description": "desc",
		"image":       "https://example.com/i.png",
	}
	b, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	art := resp["artwork"].(map[string]any)

	require.Equal(t, "Valid Title", art["title"])
	require.Equal(t, "Ártist", art["artist"]) // normalized NFC
	require.Equal(t, "desc", art["description"])
	require.Equal(t, "https://example.com/i.png", art["image"])
}
