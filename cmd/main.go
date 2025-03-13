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

	// Apply the AuthMiddleware to all routes in the group
	galleryGroup := r.Group("/gallery")
	galleryGroup.Use(middleware.AuthMiddleware("Gallery", clientID))

	// Protected routes
	galleryGroup.POST("/artworks", controllers.ArtworkCreate)
	galleryGroup.PUT("/artworks/:id", controllers.ArtworkUpdate)
	galleryGroup.DELETE("/artworks/:id", controllers.ArtworkDelete)
	galleryGroup.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"test": "test",
		})
	})

	regularGroup := r.Group("/")
	regularGroup.Use(middleware.AuthMiddleware("Regular", clientID))
	regularGroup.GET("/artworks", controllers.ArtworkIndex)
	regularGroup.GET("/artworks/:id", controllers.ArtworkFind)

	r.GET("/auth/callback", controllers.AuthCallback)

	// Run the server
	r.Run()
}
