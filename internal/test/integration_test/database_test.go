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
	env "github.com/paumarro/apollo-be/pkg"
	"github.com/stretchr/testify/assert"
)

var testArtworkID string

func TestMain(m *testing.M) {
	// Connect to the test database
	env.LoadEnvVariables("../../../.env")
	initializers.ConnectToDB()

	// Run migrations
	initializers.DB.AutoMigrate(&models.Artwork{})

	// Run tests
	code := m.Run()

	// Cleanup database
	initializers.DB.Exec("DROP TABLE artworks")

	os.Exit(code)
}

func setupContext(method string, body []byte, params map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	// Create a response recorder to capture the output
	w := httptest.NewRecorder()

	// Create a new Gin context
	ctx, _ := gin.CreateTestContext(w)

	// Set the HTTP method and body
	ctx.Request = &http.Request{
		Method: method,
		Header: make(http.Header),
	}
	if body != nil {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Use io.NopCloser
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

		ac := &controllers.ArtworkController{
			DB: initializers.DB, // Pass the database connection
		}
		ac.Create(ctx)

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
}
