package unit_test

import (
	"errors"
	"testing"

	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
	"github.com/paumarro/apollo-be/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestGetAllArtworks(t *testing.T) {
	t.Run("EmptyResult", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		mockRepo.On("FindAll").Return([]models.Artwork{}, nil)

		// Act
		artworks, err := service.GetAllArtworks()

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, artworks)
		mockRepo.AssertCalled(t, "FindAll")
	})

	t.Run("NonEmptyResult", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

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
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		mockRepo.On("FindAll").Return(nil, errors.New("database error"))

		// Act
		artworks, err := service.GetAllArtworks()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, artworks)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertCalled(t, "FindAll")
	})
}

func TestGetArtworkByID(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		mockRepo.On("FindByID", "invalid-id").Return(nil, errors.New("record not found"))

		// Act
		artwork, err := service.GetArtworkByID("invalid-id")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, artwork)
		assert.Equal(t, "record not found", err.Error())
		mockRepo.AssertCalled(t, "FindByID", "invalid-id")
	})

	t.Run("Found", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		expectedArtwork := &models.Artwork{ID: 1, Title: "Artwork 1", Artist: "Artist 1"}

		mockRepo.On("FindByID", "1").Return(expectedArtwork, nil)

		// Act
		artwork, err := service.GetArtworkByID("1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedArtwork, artwork)
		mockRepo.AssertCalled(t, "FindByID", "1")
	})
}

func TestCreateArtwork(t *testing.T) {
	t.Run("InvalidData", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		invalidArtwork := &models.Artwork{Title: ""} // Title is required
		mockRepo.On("Create", invalidArtwork).Return(errors.New("validation error"))

		// Act
		err := service.CreateArtwork(invalidArtwork)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "validation error", err.Error())
		mockRepo.AssertCalled(t, "Create", invalidArtwork)
	})

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		validArtwork := &models.Artwork{Title: "Valid Title", Artist: "Valid Artist"}
		mockRepo.On("Create", validArtwork).Return(nil)

		// Act
		err := service.CreateArtwork(validArtwork)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", validArtwork)
	})
}

func TestDeleteArtwork(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		mockRepo.On("Delete", "valid-id").Return(nil)

		// Act
		err := service.DeleteArtwork("valid-id")

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", "valid-id")
	})

	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		mockRepo.On("Delete", "invalid-id").Return(errors.New("record not found"))

		// Act
		err := service.DeleteArtwork("invalid-id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "record not found", err.Error())
		mockRepo.AssertCalled(t, "Delete", "invalid-id")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mockRepo := new(repositories.MockArtworkRepository)
		service := services.NewArtworkService(mockRepo)

		mockRepo.On("Delete", "valid-id").Return(errors.New("database error"))

		// Act
		err := service.DeleteArtwork("valid-id")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertCalled(t, "Delete", "valid-id")
	})
}
