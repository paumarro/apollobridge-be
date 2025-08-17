package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/middleware"
	"github.com/stretchr/testify/require"
)

func setupFullChainRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Mimic main.go middleware order for app routes (subset): Sanitize, Validate
	r.Use(middleware.Sanitize(), middleware.Validate())
	// Use a realistic path similar to main.go
	r.POST("/gallery/artworks", func(c *gin.Context) {
		// Echo sanitized artwork back
		if v, ok := c.Get("sanitizedArtwork"); ok {
			c.JSON(http.StatusOK, gin.H{"artwork": v})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing sanitizedArtwork"})
	})
	return r
}

func TestIntegration_MiddlewareChain_Success(t *testing.T) {
	r := setupFullChainRouter()
	body := map[string]any{
		"title":       "Valid Title",
		"artist":      "Valid Artist",
		"description": "Valid description",
		"image":       "https://example.com/i.png",
	}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	art := resp["artwork"].(map[string]any)
	require.Equal(t, "Valid Title", art["title"])
	require.Equal(t, "Valid Artist", art["artist"])
	require.Equal(t, "Valid description", art["description"])
	require.Equal(t, "https://example.com/i.png", art["image"])
}

func TestIntegration_MiddlewareChain_AcceptHeaderNotAcceptable(t *testing.T) {
	r := setupFullChainRouter()
	b, _ := json.Marshal(map[string]any{
		"title":       "Valid Title",
		"artist":      "Valid Artist",
		"description": "desc",
		"image":       "https://example.com/i.png",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotAcceptable, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "not_acceptable", errObj["code"])
	require.Equal(t, "Requested representation not acceptable", errObj["message"])
}
