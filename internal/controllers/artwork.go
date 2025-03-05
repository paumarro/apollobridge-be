package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	initializers "github.com/paumarro/apollo-be/internal/art-service/initializers"
	"github.com/paumarro/apollo-be/internal/art-service/models"
)

var body struct {
	Title       string
	Artist      string
	Description string
}

func ArtworkCreate(c *gin.Context) {
	c.Bind(&body)

	artwork := models.Artwork{Title: body.Title, Artist: body.Artist, Description: body.Description}

	result := initializers.DB.Create(&artwork)

	if result.Error != nil {
		log.Printf("Error creating artwork: %v", result.Error)
		c.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"Artwork": artwork,
	})
}

func ArtworkIndex(c *gin.Context) {
	var artworks []models.Artwork
	result := initializers.DB.Find(&artworks)

	if result.Error != nil {
		log.Fatal(result.Error)
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"artworks": artworks,
	})
}

func ArtworkFind(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	result := initializers.DB.First(&artwork, id)

	if result.Error != nil {
		log.Fatal(result.Error)
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"artworks": artwork,
	})
}

func ArtworkUpdate(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork
	initializers.DB.First(&artwork, id)

	initializers.DB.Model(&artwork).Updates(models.Artwork{
		Title:       body.Title,
		Artist:      body.Artist,
		Description: body.Description,
	})

	c.JSON(200, gin.H{
		"artworks": artwork,
	})
}

func ArtworkDelete(c *gin.Context) {
	id := c.Param("id")

	var artwork models.Artwork

	initializers.DB.Delete(&artwork, id)

	c.JSON(200, gin.H{
		"message": "Artwork successfuly deleted",
	})
}
