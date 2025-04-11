package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/dto"
	initializers "github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/models"
	"gorm.io/gorm"
)

// Define the struct with validation tags
func respondWithError(c *gin.Context, code int, message string, details interface{}) {
	log.Printf("Error: %s, Details: %v", message, details) // Log the full error for debugging
	c.JSON(code, gin.H{"error": message})                  // Send a user-friendly error message
}

func ArtworkCreate(c *gin.Context) {
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		c.JSON(500, gin.H{"error": "Failed to retrieve sanitized input"})
		return
	}

	req := sanitizedArtwork.(dto.ArtworkRequest)

	artwork := models.Artwork{
		Title:       req.Title,
		Artist:      req.Artist,
		Description: req.Description,
		Image:       req.Image,
	}

	if result := initializers.DB.Create(&artwork); result.Error != nil {
		log.Printf("Error creating artwork: %v", result.Error)
		respondWithError(c, http.StatusInternalServerError, "Failed to create artwork", nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"artwork": artwork})
}

func ArtworkIndex(c *gin.Context) {
	var artworks []models.Artwork
	if result := initializers.DB.Find(&artworks); result.Error != nil {
		log.Printf("Error fetching artworks: %v", result.Error)
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch artworks", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artworks": artworks})
}

func ArtworkFind(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	result := initializers.DB.First(&artwork, id)

	if result.Error != nil {
		log.Printf("Error finding artwork: %v", result.Error)
		respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

func ArtworkUpdate(c *gin.Context) {
	// Retrieve the artwork ID from the path parameters
	id := c.Param("id")

	// Fetch the existing artwork from the database
	var artwork models.Artwork
	if err := initializers.DB.First(&artwork, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			log.Printf("Error finding artwork: %v", err)
			respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		}
		return
	}

	// Retrieve the sanitized and validated input from the middleware
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sanitized input"})
		return
	}
	req := sanitizedArtwork.(dto.ArtworkRequest)

	// Update the artwork fields with the sanitized input
	artwork.Title = req.Title
	artwork.Artist = req.Artist
	artwork.Description = req.Description
	artwork.Image = req.Image

	// Save the updated artwork to the database
	if err := initializers.DB.Save(&artwork).Error; err != nil {
		log.Printf("Error updating artwork: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to update artwork", nil)
		return
	}

	// Respond with the updated artwork
	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

func ArtworkDelete(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	if err := initializers.DB.First(&artwork, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			log.Printf("Error finding artwork: %v", err)
			respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		}
		return
	}

	if err := initializers.DB.Delete(&artwork).Error; err != nil {
		log.Printf("Error deleting artwork: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to delete artwork", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artwork successfully deleted"})
}
