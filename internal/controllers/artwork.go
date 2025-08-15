package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/dto"
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/services"
)

// ArtworkController handles artwork-related operations
type ArtworkController struct {
	ArtworkService *services.ArtworkService
}

// NewArtworkController creates a new instance of ArtworkController
func NewArtworkController(service *services.ArtworkService) *ArtworkController {
	return &ArtworkController{ArtworkService: service}
}

// respondWithError is a helper function to send error responses
func respondWithError(c *gin.Context, code int, message string, details interface{}) {
	log.Printf("Error: %s, Details: %v", message, details)
	c.JSON(code, gin.H{"error": message})
}

// ArtworkCreate handles the creation of a new artwork
func (ac *ArtworkController) Create(c *gin.Context) {
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve sanitized input", nil)
		return
	}

	req, ok := sanitizedArtwork.(dto.ArtworkRequest)
	if !ok {
		respondWithError(c, http.StatusInternalServerError, "Invalid sanitized input", nil)
		return
	}
	artwork := models.Artwork{
		Title:       req.Title,
		Artist:      req.Artist,
		Description: req.Description,
		Image:       req.Image,
	}

	if err := ac.ArtworkService.CreateArtwork(&artwork); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create artwork", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"artwork": artwork})
}

func (ac *ArtworkController) Index(c *gin.Context) {
	artworks, err := ac.ArtworkService.GetAllArtworks()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch artworks", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artworks": artworks})
}

func (ac *ArtworkController) Find(c *gin.Context) {
	id := c.Param("id")

	artwork, err := ac.ArtworkService.GetArtworkByID(id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

// ArtworkUpdate handles updating an artwork
func (ac *ArtworkController) Update(c *gin.Context) {
	id := c.Param("id")

	// Fetch the artwork by ID
	artwork, err := ac.ArtworkService.GetArtworkByID(id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			log.Printf("Error finding artwork: %v", err)
			respondWithError(c, http.StatusInternalServerError, "Failed to find artwork", nil)
		}
		return
	}

	// Retrieve sanitized input
	sanitizedArtwork, exists := c.Get("sanitizedArtwork")
	if !exists {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve sanitized input", nil)
		return
	}

	// Update the artwork fields
	req := sanitizedArtwork.(dto.ArtworkRequest)
	artwork.Title = req.Title
	artwork.Artist = req.Artist
	artwork.Description = req.Description
	artwork.Image = req.Image

	// Attempt to update the artwork
	if err := ac.ArtworkService.UpdateArtwork(artwork); err != nil {
		log.Printf("Error updating artwork: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to update artwork", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

func (ac *ArtworkController) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := ac.ArtworkService.DeleteArtwork(id); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, "Artwork not found", nil)
		} else {
			respondWithError(c, http.StatusInternalServerError, "Failed to delete artwork", err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artwork successfully deleted"})
}
