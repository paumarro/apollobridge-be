package unit_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/controllers"
	"github.com/paumarro/apollo-be/internal/dto"
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
	"github.com/paumarro/apollo-be/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SetupMockController initializes a mock repository and returns an ArtworkController instance.
func setupMockController() (*repositories.MockArtworkRepository, *controllers.ArtworkController) {
	mockRepo := &repositories.MockArtworkRepository{}
	artworkService := services.NewArtworkService(mockRepo)
	artworkController := controllers.NewArtworkController(artworkService)
	return mockRepo, artworkController
}

func TestArtworkCreate(t *testing.T) {
	t.Run("Successful Creation", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		// Configure the mock to return no error
		mockRepo.On("Create", mock.AnythingOfType("*models.Artwork")).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the sanitizedArtwork in the Gin context
		c.Set("sanitizedArtwork", dto.ArtworkRequest{
			Title:       "Test Artwork",
			Artist:      "Test Artist",
			Description: "Test Description",
			Image:       "http://test.com/image.jpg",
		})

		// Call the controller
		ac.Create(c)

		// Assert that the response is 201 Created
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Test Artwork")

		// Verify the mock was called
		mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("*models.Artwork"))
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		// Configure the mock to return a database error
		mockRepo.On("Create", mock.AnythingOfType("*models.Artwork")).Run(func(args mock.Arguments) {
			fmt.Println("Mock Create called with:", args)
		}).Return(errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the sanitizedArtwork in the Gin context
		c.Set("sanitizedArtwork", dto.ArtworkRequest{
			Title:       "Test Artwork",
			Artist:      "Test Artist",
			Description: "Test Description",
			Image:       "http://test.com/image.jpg",
		})

		// Call the controller
		ac.Create(c)

		// Assert that the response is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create artwork")

		// Verify the mock was called
		mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("*models.Artwork"))
	})

	t.Run("Missing Sanitized Artwork in Context", func(t *testing.T) {
		_, ac := setupMockController()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Do not set "sanitizedArtwork" in the context

		ac.Create(c)

		// Assert that the response is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve sanitized input")
	})
}

func TestArtworkIndex(t *testing.T) {
	t.Run("Successful Fetch with Artworks", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("FindAll").Return([]models.Artwork{
			{Title: "Artwork 1", Artist: "Artist 1"},
			{Title: "Artwork 2", Artist: "Artist 2"},
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ac.Index(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork 1")
		assert.Contains(t, w.Body.String(), "Artwork 2")

		mockRepo.AssertCalled(t, "FindAll")
	})

	t.Run("Empty Database", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("FindAll").Return([]models.Artwork{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ac.Index(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"artworks":[]`)

		mockRepo.AssertCalled(t, "FindAll")
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("FindAll").Return([]models.Artwork{}, errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ac.Index(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch artworks")

		mockRepo.AssertCalled(t, "FindAll")
	})
}

func TestArtworkFind(t *testing.T) {
	t.Run("Successful Retrieval", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("FindByID", "1").Return(&models.Artwork{
			Title:  "Artwork 1",
			Artist: "Artist 1",
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		ac.Find(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork 1")

		mockRepo.AssertCalled(t, "FindByID", "1")
	})

	t.Run("Artwork Not Found", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("FindByID", "1").Return(nil, services.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		ac.Find(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork not found")

		mockRepo.AssertCalled(t, "FindByID", "1")
	})
}

func TestArtworkUpdate(t *testing.T) {
	t.Run("Successful Update", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		// Mock the repository to return an existing artwork
		mockRepo.On("FindByID", "1").Return(&models.Artwork{
			ID:          1,
			Title:       "Old Title",
			Artist:      "Old Artist",
			Description: "Old Description",
			Image:       "http://oldimage.com",
		}, nil)

		// Mock the repository to update the artwork successfully
		mockRepo.On("Update", mock.AnythingOfType("*models.Artwork")).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the ID parameter
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Set the sanitizedArtwork in the Gin context
		c.Set("sanitizedArtwork", dto.ArtworkRequest{
			Title:       "New Title",
			Artist:      "New Artist",
			Description: "New Description",
			Image:       "http://newimage.com",
		})

		// Call the controller
		ac.Update(c)

		// Assert that the response is 200 OK
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "New Title")
		assert.Contains(t, w.Body.String(), "New Artist")

		// Verify the mock was called
		mockRepo.AssertCalled(t, "FindByID", "1")
		mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("*models.Artwork"))
	})

	t.Run("Artwork Not Found", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		// Mock the repository to return not found
		mockRepo.On("FindByID", "1").Return(nil, services.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the ID parameter
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Call the controller
		ac.Update(c)

		// Assert that the response is 404 Not Found
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork not found")

		// Verify the mock was called
		mockRepo.AssertCalled(t, "FindByID", "1")
	})

	t.Run("Missing Sanitized Input", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		// Mock the repository to return an existing artwork
		mockRepo.On("FindByID", "1").Return(&models.Artwork{
			ID:          1,
			Title:       "Old Title",
			Artist:      "Old Artist",
			Description: "Old Description",
			Image:       "http://oldimage.com",
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the ID parameter
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Do not set sanitizedArtwork in the context
		ac.Update(c)

		// Assert that the response is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve sanitized input")

		// Verify the mock was called
		mockRepo.AssertCalled(t, "FindByID", "1")
	})

	t.Run("Database Error During Update", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		// Mock the repository to return an existing artwork
		mockRepo.On("FindByID", "1").Return(&models.Artwork{
			ID:          1,
			Title:       "Old Title",
			Artist:      "Old Artist",
			Description: "Old Description",
			Image:       "http://oldimage.com",
		}, nil)

		// Mock the repository to fail during the update
		mockRepo.On("Update", mock.AnythingOfType("*models.Artwork")).Return(errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the ID parameter
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Set the sanitizedArtwork in the Gin context
		c.Set("sanitizedArtwork", dto.ArtworkRequest{
			Title:       "New Title",
			Artist:      "New Artist",
			Description: "New Description",
			Image:       "http://newimage.com",
		})

		// Call the controller
		ac.Update(c)

		// Assert that the response is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update artwork")

		// Verify the mock was called
		mockRepo.AssertCalled(t, "FindByID", "1")
		mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("*models.Artwork"))
	})
}

func TestArtworkDelete(t *testing.T) {
	t.Run("Successful Deletion", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("Delete", "1").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		ac.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork successfully deleted")

		mockRepo.AssertCalled(t, "Delete", "1")
	})

	t.Run("Artwork Not Found", func(t *testing.T) {
		mockRepo, ac := setupMockController() // Fresh mockRepo for this subtest

		mockRepo.On("Delete", "1").Return(services.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		ac.Delete(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork not found")

		mockRepo.AssertCalled(t, "Delete", "1")
	})
}
