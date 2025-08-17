package unit_test

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"gorm.io/gorm"

	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
	"github.com/paumarro/apollo-be/internal/services"
	"github.com/stretchr/testify/assert"
)

func setUpMockServiceWithLogger() (*repositories.MockArtworkRepository, *services.ArtworkService, *bytes.Buffer) {
	mockRepo := &repositories.MockArtworkRepository{}
	logBuffer := new(bytes.Buffer) // Buffer to capture logs
	log.SetOutput(logBuffer)       // Redirect logs to the buffer
	service := services.NewArtworkService(mockRepo)
	return mockRepo, service, logBuffer
}

func TestGetAllArtworks(t *testing.T) {
	t.Run("EmptyResult", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		mockRepo.On("FindAll").Return([]models.Artwork{}, nil)

		// Act
		artworks, err := service.GetAllArtworks()

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, artworks)
		mockRepo.AssertCalled(t, "FindAll")
		assert.Contains(t, logBuffer.String(), "Fetching all artworks")
	})

	t.Run("NonEmptyResult", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		expectedArtworks := []models.Artwork{
			{ID: 1, Title: "Artwork 1", Artist: "Artist 1"},
			{ID: 2, Title: "Artwork 2", Artist: "Artist 2"},
		}

		mockRepo.On("FindAll").Return(expectedArtworks, nil)

		// Act
		artworks, err := service.GetAllArtworks()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedArtworks, artworks)
		mockRepo.AssertCalled(t, "FindAll")
		assert.Contains(t, logBuffer.String(), "Fetching all artworks")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		mockRepo.On("FindAll").Return(nil, errors.New("database error"))

		// Act
		artworks, err := service.GetAllArtworks()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, artworks)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertCalled(t, "FindAll")
		assert.Contains(t, logBuffer.String(), "Fetching all artworks")
	})
}

func TestGetArtworkByID(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		mockRepo.On("FindByID", "invalid-id").Return(nil, gorm.ErrRecordNotFound)

		// Act
		artwork, err := service.GetArtworkByID("invalid-id")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, artwork)
		assert.True(t, errors.Is(err, services.ErrNotFound), "expected ErrNotFound error")
		mockRepo.AssertCalled(t, "FindByID", "invalid-id")
		assert.Contains(t, logBuffer.String(), "Fetching artwork with ID: invalid-id")
	})

	t.Run("Found", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		expectedArtwork := &models.Artwork{ID: 1, Title: "Artwork 1", Artist: "Artist 1"}

		mockRepo.On("FindByID", "1").Return(expectedArtwork, nil)

		// Act
		artwork, err := service.GetArtworkByID("1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedArtwork, artwork)
		mockRepo.AssertCalled(t, "FindByID", "1")
		assert.Contains(t, logBuffer.String(), "Fetching artwork with ID: 1")
	})
}

func TestCreateArtwork(t *testing.T) {
	t.Run("InvalidData", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		invalidArtwork := &models.Artwork{Title: "", Artist: "Someone"} // Title is required (validated by repo/db)

		mockRepo.On("FindAll").Return([]models.Artwork{}, nil)
		mockRepo.On("Create", invalidArtwork).Return(errors.New("validation error"))

		// Act
		err := service.CreateArtwork(invalidArtwork)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation error")
		mockRepo.AssertCalled(t, "FindAll")
		mockRepo.AssertCalled(t, "Create", invalidArtwork)
		assert.Contains(t, logBuffer.String(), "Creating artwork:")
		assert.Contains(t, logBuffer.String(), "Error in CreateArtwork: validation error")
	})

	t.Run("DuplicateArtwork", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		dup := &models.Artwork{Title: "Same", Artist: "Artist"}
		// Service checks duplicates using FindAll (title + artist)
		mockRepo.On("FindAll").Return([]models.Artwork{
			{ID: 10, Title: "Same", Artist: "Artist"},
		}, nil)
		// Create should NOT be called

		// Act
		err := service.CreateArtwork(dup)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, services.ErrDuplicate))
		mockRepo.AssertCalled(t, "FindAll")
		mockRepo.AssertNotCalled(t, "Create")
		assert.Contains(t, logBuffer.String(), "Creating artwork:")
	})

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		validArtwork := &models.Artwork{Title: "Valid Title", Artist: "Valid Artist"}
		// No duplicates
		mockRepo.On("FindAll").Return([]models.Artwork{}, nil)
		mockRepo.On("Create", validArtwork).Return(nil)

		// Act
		err := service.CreateArtwork(validArtwork)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "FindAll")
		mockRepo.AssertCalled(t, "Create", validArtwork)
		assert.Contains(t, logBuffer.String(), "Creating artwork:")
	})
}

func TestDeleteArtwork(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		mockRepo.On("Delete", "valid-id").Return(nil)

		// Act
		err := service.DeleteArtwork("valid-id")

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", "valid-id")
		assert.Contains(t, logBuffer.String(), "Deleting artwork with ID: valid-id")
	})

	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		mockRepo.On("Delete", "invalid-id").Return(gorm.ErrRecordNotFound)

		// Act
		err := service.DeleteArtwork("invalid-id")

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, services.ErrNotFound), "expected ErrNotFound error")
		mockRepo.AssertCalled(t, "Delete", "invalid-id")
		assert.Contains(t, logBuffer.String(), "Deleting artwork with ID: invalid-id")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mockRepo, service, logBuffer := setUpMockServiceWithLogger()

		mockRepo.On("Delete", "valid-id").Return(errors.New("database error"))

		// Act
		err := service.DeleteArtwork("valid-id")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertCalled(t, "Delete", "valid-id")
		assert.Contains(t, logBuffer.String(), "Deleting artwork with ID: valid-id")
	})
}
