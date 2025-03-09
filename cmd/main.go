package main

import (
	"log"

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

	// Apply the AuthMiddleware to all routes in the group
	secured := r.Group("/")
	secured.Use(middleware.AuthMiddleware())

	// Protected routes
	secured.POST("/artworks", controllers.ArtworkCreate)
	secured.GET("/artworks", controllers.ArtworkIndex)
	secured.GET("/artworks/:id", controllers.ArtworkFind)
	secured.PUT("/artworks/:id", controllers.ArtworkUpdate)
	secured.DELETE("/artworks/:id", controllers.ArtworkDelete)

	r.GET("/auth/callback", controllers.AuthCallback)

	// Run the server
	r.Run()
}
