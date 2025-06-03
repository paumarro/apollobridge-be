package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/controllers"
	"github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/middleware"
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/services"
	env "github.com/paumarro/apollo-be/pkg"
)

func init() {
	env.LoadEnvVariables(".env")
	initializers.ConnectToDB()
	if err := initializers.DB.AutoMigrate(&models.Artwork{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
	r := gin.Default()

	// clientID := os.Getenv("KEYCLOAK_CLIENT_ID")

	r.Use(middleware.RateLimit())
	r.Use(middleware.SecurityHeaders())

	// Instantiate the service and controller
	artworkService := services.NewArtworkService(initializers.DB)
	artworkController := controllers.NewArtworkController(artworkService)

	galleryGroup := r.Group("/gallery")
	// galleryGroup.Use(middleware.Auth("Gallery", clientID))
	galleryGroup.Use(middleware.Sanitize())
	galleryGroup.Use(middleware.Validate())

	galleryGroup.POST("/artworks", artworkController.Create)
	galleryGroup.PUT("/artworks/:id", artworkController.Update)
	galleryGroup.DELETE("/artworks/:id", artworkController.Delete)

	regularGroup := r.Group("/")
	// regularGroup.Use(middleware.Auth("Regular", clientID))
	regularGroup.Use(middleware.Sanitize())
	regularGroup.Use(middleware.Validate())

	regularGroup.GET("/artworks", artworkController.Index)
	regularGroup.GET("/artworks/:id", artworkController.Find)

	r.GET("/auth/callback", middleware.Sanitize(), middleware.Validate(), controllers.AuthCallback)

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
