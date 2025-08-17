package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/dto"
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/services"
)

func parseIDParam(c *gin.Context) (uint64, bool) {
	idStr := c.Param("id")
	if idStr == "" {
		respondWithError(c, http.StatusBadRequest, "Invalid or missing ID parameter", nil)
		return 0, false
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format, must be a positive integer", nil)
		return 0, false
	}
	return id, true
}

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
		if strings.Contains(err.Error(), "already exists") {
			// Map duplicate error to 409 Conflict
			respondWithError(c, http.StatusConflict, "Artwork already exists", err.Error())
		} else {
			// Handle other errors as 500 Internal Server Error
			respondWithError(c, http.StatusInternalServerError, "Failed to create artwork", err.Error())
		}
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

// Find
func (ac *ArtworkController) Find(c *gin.Context) {
	_, ok := parseIDParam(c)
	if !ok {
		return
	}

	id := c.Param("id") // still pass the string if your service expects string today
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
// Update
func (ac *ArtworkController) Update(c *gin.Context) {
	_, ok := parseIDParam(c)
	if !ok {
		return
	}

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

	artwork.Title = req.Title
	artwork.Artist = req.Artist
	artwork.Description = req.Description
	artwork.Image = req.Image

	if err := ac.ArtworkService.UpdateArtwork(artwork); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update artwork", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"artwork": artwork})
}

// Delete
func (ac *ArtworkController) Delete(c *gin.Context) {
	_, ok := parseIDParam(c)
	if !ok {
		return
	}

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
