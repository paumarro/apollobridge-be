package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/controllers"
	"github.com/paumarro/apollo-be/internal/dto"
	"github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
	"github.com/paumarro/apollo-be/internal/services"
	env "github.com/paumarro/apollo-be/pkg"
	"github.com/stretchr/testify/assert"
)

var testArtworkID string

func TestMain(m *testing.M) {
	// Load environment variables
	env.LoadEnvVariables("../../../.env")

	// Connect to the database
	initializers.ConnectToDB()

	// Run migrations
	initializers.DB.AutoMigrate(&models.Artwork{})

	// Run the tests
	code := m.Run()

	// Cleanup database
	initializers.DB.Exec("DROP TABLE artworks")

	os.Exit(code)
}

func setupContext(method string, body []byte, params map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a new Gin context
	ctx, _ := gin.CreateTestContext(w)

	// Set the HTTP method and body
	ctx.Request = &http.Request{
		Method: method,
		Header: make(http.Header),
	}
	if body != nil {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")
	}

	// Set URL parameters
	for key, value := range params {
		ctx.Params = append(ctx.Params, gin.Param{Key: key, Value: value})
	}

	return ctx, w
}

func TestArtworkIntegration(t *testing.T) {
	log.Printf("initializers.DB: %+v", initializers.DB)

	repo := repositories.NewGormArtworkRepository(initializers.DB)
	service := services.NewArtworkService(repo)
	controller := controllers.NewArtworkController(service)

	// Test creating an artwork
	t.Run("Create Artwork", func(t *testing.T) {
		artwork := dto.ArtworkRequest{
			Title:       "Mona Lisa",
			Artist:      "Leonardo da Vinci",
			Description: "A portrait of a woman",
			Image:       "monalisa.jpg",
		}

		body, _ := json.Marshal(artwork)
		ctx, w := setupContext("POST", body, nil)

		ctx.Set("sanitizedArtwork", artwork)
		controller.Create(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		artworkResponse := response["artwork"].(map[string]interface{})
		assert.Equal(t, artwork.Title, artworkResponse["title"])
		assert.Equal(t, artwork.Artist, artworkResponse["artist"])
		assert.Equal(t, artwork.Description, artworkResponse["description"])
		assert.Equal(t, artwork.Image, artworkResponse["image"])

		testArtworkID = fmt.Sprintf("%v", artworkResponse["id"])
	})

	// Test edge case: Duplicate Artwork
	t.Run("Create Duplicate Artwork", func(t *testing.T) {
		artwork := dto.ArtworkRequest{
			Title:       "Mona Lisa", // Same title as the previous test
			Artist:      "Leonardo da Vinci",
			Description: "Duplicate test",
			Image:       "duplicate.jpg",
		}

		body, _ := json.Marshal(artwork)
		ctx, w := setupContext("POST", body, nil)

		ctx.Set("sanitizedArtwork", artwork)
		controller.Create(ctx)

		// Expect a conflict or bad request status code
		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response["error"], "Artwork already exists")
	})

	// Test fetching all artworks
	t.Run("Fetch All Artworks", func(t *testing.T) {
		ctx, w := setupContext("GET", nil, nil)
		controller.Index(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		artworks := response["artworks"].([]interface{})
		assert.GreaterOrEqual(t, len(artworks), 1)
		artwork := artworks[0].(map[string]interface{})
		assert.Equal(t, "Mona Lisa", artwork["title"])
		assert.Equal(t, "Leonardo da Vinci", artwork["artist"])
	})

	// Test edge case: Fetch Non-Existent Artwork
	t.Run("Fetch Non-Existent Artwork", func(t *testing.T) {
		ctx, w := setupContext("GET", nil, map[string]string{"id": "2222222222222222222"})
		controller.Find(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Artwork not found", response["error"])
	})

	t.Run("Fetch Artwork with Invalid ID", func(t *testing.T) {
		ctx, w := setupContext("GET", nil, map[string]string{"id": "invalid-id"})
		controller.Find(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Invalid ID format, must be a positive integer", response["error"])
	})

	// Test updating an artwork
	t.Run("Update Artwork", func(t *testing.T) {
		updatedArtwork := dto.ArtworkRequest{
			Title:       "Mona Lisa Updated",
			Artist:      "Leonardo da Vinci",
			Description: "An updated description",
			Image:       "monalisa_updated.jpg",
		}

		body, _ := json.Marshal(updatedArtwork)
		ctx, w := setupContext("PUT", body, map[string]string{"id": testArtworkID})

		ctx.Set("sanitizedArtwork", updatedArtwork)
		controller.Update(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		artwork := response["artwork"].(map[string]interface{})
		assert.Equal(t, updatedArtwork.Title, artwork["title"])
		assert.Equal(t, updatedArtwork.Description, artwork["description"])
		assert.Equal(t, updatedArtwork.Image, artwork["image"])
	})

	// Test edge case: Update Non-Existent Artwork
	t.Run("Update Non-Existent Artwork", func(t *testing.T) {
		updatedArtwork := dto.ArtworkRequest{
			Title:       "Non-Existent",
			Artist:      "Unknown",
			Description: "This artwork does not exist",
			Image:       "nonexistent.jpg",
		}

		body, _ := json.Marshal(updatedArtwork)
		ctx, w := setupContext("PUT", body, map[string]string{"id": "222222222222222222"})

		ctx.Set("sanitizedArtwork", updatedArtwork)
		controller.Update(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Artwork not found", response["error"])
	})

	// Test deleting an artwork
	t.Run("Delete Artwork", func(t *testing.T) {
		ctx, w := setupContext("DELETE", nil, map[string]string{"id": testArtworkID})
		controller.Delete(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Artwork successfully deleted", response["message"])
	})

	// Test edge case: Delete Non-Existent Artwork
	t.Run("Delete Non-Existent Artwork", func(t *testing.T) {
		ctx, w := setupContext("DELETE", nil, map[string]string{"id": "222222222222222222"})
		controller.Delete(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Artwork not found", response["error"])
	})
}
