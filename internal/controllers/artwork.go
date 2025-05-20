package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/dto"
	"github.com/paumarro/apollo-be/internal/models"
	"gorm.io/gorm"
)

// ArtworkController handles artwork-related operations
type ArtworkController struct {
	DB *gorm.DB
}

// NewArtworkController creates a new instance of ArtworkController
func NewArtworkController(db *gorm.DB) *ArtworkController {
	return &ArtworkController{DB: db}
}

// RespondWithError is a helper function to send error responses
func respondWithError(c *gin.Context, code int, message string, details interface{}) {
	log.Printf("Error: %s, Details: %v", message, details)
	c.JSON(code, gin.H{"error": message})
}

// ArtworkCreate handles the creation of a new artwork
func (ac *ArtworkController) ArtworkCreate(c *gin.Context) {
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		log.Println("Sanitized artwork not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sanitized input"})
		return
	}

	req := sanitizedArtwork.(dto.ArtworkRequest)
	log.Printf("Received artwork request: %+v", req)

	artwork := models.Artwork{
		Title:       req.Title,
		Artist:      req.Artist,
		Description: req.Description,
		Image:       req.Image,
	}

	if result := ac.DB.Create(&artwork); result.Error != nil {
		log.Printf("Error creating artwork: %v", result.Error)
		respondWithError(c, http.StatusInternalServerError, "Failed to create artwork", nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"artwork": artwork})
}

// ArtworkIndex handles fetching all artworks
func (ac *ArtworkController) ArtworkIndex(c *gin.Context) {
	var artworks []models.Artwork
	if result := ac.DB.Find(&artworks); result.Error != nil {
		log.Printf("Error fetching artworks: %v", result.Error)
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch artworks", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artworks": artworks})
}

// ArtworkFind handles fetching a single artwork by ID
func (ac *ArtworkController) ArtworkFind(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	result := ac.DB.First(&artwork, id)

	if result.Error != nil {
		log.Printf("Error finding artwork: %v", result.Error)
		respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

// ArtworkUpdate handles updating an artwork
func (ac *ArtworkController) ArtworkUpdate(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	if err := ac.DB.First(&artwork, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			log.Printf("Error finding artwork: %v", err)
			respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		}
		return
	}

	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sanitized input"})
		return
	}
	req := sanitizedArtwork.(dto.ArtworkRequest)

	artwork.Title = req.Title
	artwork.Artist = req.Artist
	artwork.Description = req.Description
	artwork.Image = req.Image

	if err := ac.DB.Save(&artwork).Error; err != nil {
		log.Printf("Error updating artwork: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to update artwork", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

// ArtworkDelete handles deleting an artwork
func (ac *ArtworkController) ArtworkDelete(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	if err := ac.DB.First(&artwork, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			log.Printf("Error finding artwork: %v", err)
			respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		}
		return
	}

	if err := ac.DB.Delete(&artwork).Error; err != nil {
		log.Printf("Error deleting artwork: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to delete artwork", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artwork successfully deleted"})
}
