package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/controllers"
	"github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/middleware"
	"github.com/paumarro/apollo-be/internal/models"
	env "github.com/paumarro/apollo-be/pkg"
)

func init() {
	// Load environment variables
	env.LoadEnvVariables()

	// Connect to the database
	initializers.ConnectToDB()

	// Run database migrations
	if err := initializers.DB.AutoMigrate(&models.Artwork{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Retrieve the Keycloak client ID from the environment
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")

	// Initialize the controllers with the database instance
	artworkController := controllers.NewArtworkController(initializers.DB)

	// Apply global middleware
	r.Use(middleware.RateLimit())
	r.Use(middleware.SecurityHeaders())

	// Gallery group routes with specific middleware
	galleryGroup := r.Group("/gallery")
	galleryGroup.Use(middleware.Auth("Gallery", clientID))
	galleryGroup.Use(middleware.Sanitize())
	galleryGroup.Use(middleware.Validate())

	galleryGroup.POST("/artworks", artworkController.ArtworkCreate)
	galleryGroup.PUT("/artworks/:id", artworkController.ArtworkUpdate)
	galleryGroup.DELETE("/artworks/:id", artworkController.ArtworkDelete)

	// Regular group routes with specific middleware
	regularGroup := r.Group("/")
	regularGroup.Use(middleware.Auth("Regular", clientID))
	regularGroup.Use(middleware.Sanitize())
	regularGroup.Use(middleware.Validate())

	regularGroup.GET("/artworks", artworkController.ArtworkIndex)
	regularGroup.GET("/artworks/:id", artworkController.ArtworkFind)

	// Authentication callback route
	r.GET("/auth/callback", middleware.Sanitize(), middleware.Validate(), controllers.AuthCallback)

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
