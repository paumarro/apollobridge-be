package unit_test

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

func setupValidateRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Full chain: Sanitize then Validate, like main.go
	r.Use(middleware.Sanitize(), middleware.Validate())
	r.POST("/gallery/artworks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	// Add a GET route to test Accept header handling (Sanitize is inert on GET)
	r.GET("/artworks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func validBody() []byte {
	b, _ := json.Marshal(map[string]any{
		"title":       "Valid Title",
		"artist":      "Valid Artist",
		"description": "A nice description",
		"image":       "https://example.com/i.png",
	})
	return b
}

func TestValidate_AcceptHeader_NotAcceptable_GET(t *testing.T) {
	r := setupValidateRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/artworks", nil)
	req.Header.Set("Accept", "text/plain")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotAcceptable, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "not_acceptable", errObj["code"])
	require.Equal(t, "Requested representation not acceptable", errObj["message"])
}

func TestValidate_AcceptHeader_AllowsJSONish_GET(t *testing.T) {
	r := setupValidateRouter()
	for _, accept := range []string{"application/json", "application/vnd.api+json", "*/*"} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/artworks", nil)
		req.Header.Set("Accept", accept)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	}
}

func TestValidate_QueryParams_Failures(t *testing.T) {
	r := setupValidateRouter()

	// invalid state
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks?state=bad", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "validation_error", errObj["code"])

	// invalid session_state
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/gallery/artworks?session_state=not-a-uuid", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// invalid iss (non-https)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/gallery/artworks?iss=http://issuer.example", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// invalid code (short)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/gallery/artworks?code=short", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidate_DTOValidation_Failures(t *testing.T) {
	r := setupValidateRouter()

	body := map[string]any{
		"title":       "",                // required, min
		"artist":      "ab",              // min 3
		"description": "",                // required
		"image":       "https://e.com/i", // will pass URL format, but keep a separate test for https enforcement too
	}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "validation_error", errObj["code"])
	require.Equal(t, "Request validation failed", errObj["message"])
	require.Contains(t, errObj, "details")
	details := errObj["details"].(map[string]any)
	require.NotEmpty(t, details)
	// At least these fields should be flagged
	require.Contains(t, details, "title")
	require.Contains(t, details, "artist")
	require.Contains(t, details, "description")
}

func TestValidate_ImageMustBeHTTPS(t *testing.T) {
	r := setupValidateRouter()

	body := map[string]any{
		"title":       "TTT",
		"artist":      "AAA",
		"description": "desc",
		"image":       "http://example.com/i.png", // not https
	}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "validation_error", errObj["code"])
	details := errObj["details"].(map[string]any)
	require.Equal(t, "image must be a valid https URL", details["image"])
}

func TestValidate_ControlChars(t *testing.T) {
	r := setupValidateRouter()

	body := map[string]any{
		"title":       "Good\x00Title", // contains NUL
		"artist":      "Valid Artist",
		"description": "ok\ndesc", // allowed newlines
		"image":       "https://example.com/i.png",
	}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/gallery/artworks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	require.Equal(t, "validation_error", errObj["code"])
	details := errObj["details"].(map[string]any)
	require.Equal(t, "title contains invalid control characters", details["title"])
}
