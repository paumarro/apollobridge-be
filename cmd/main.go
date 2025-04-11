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
	env.LoadEnvVariables()
	initializers.ConnectToDB()
	if err := initializers.DB.AutoMigrate(&models.Artwork{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
	r := gin.Default()

	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")

	r.Use(middleware.RateLimit())

	galleryGroup := r.Group("/gallery")
	galleryGroup.Use(middleware.Auth("Gallery", clientID))
	galleryGroup.Use(middleware.Sanitize())
	galleryGroup.Use(middleware.Validate())

	galleryGroup.POST("/artworks", controllers.ArtworkCreate)
	galleryGroup.PUT("/artworks/:id", controllers.ArtworkUpdate)
	galleryGroup.DELETE("/artworks/:id", controllers.ArtworkDelete)

	regularGroup := r.Group("/")
	regularGroup.Use(middleware.Auth("Regular", clientID))
	regularGroup.Use(middleware.Sanitize())
	regularGroup.Use(middleware.Validate())

	regularGroup.GET("/artworks", controllers.ArtworkIndex)
	regularGroup.GET("/artworks/:id", controllers.ArtworkFind)

	r.GET("/auth/callback", controllers.AuthCallback)

	r.Run()
}
