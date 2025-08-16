package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/controllers"
	"github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/middleware"
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
	"github.com/paumarro/apollo-be/internal/services"
	env "github.com/paumarro/apollo-be/pkg"
)

func init() {
	env.LoadEnvVariables(".env")
	initializers.ConnectToDB()
	log.Println("Starting database migration...")
	if err := initializers.DB.AutoMigrate(&models.Artwork{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully.")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// clientID := os.Getenv("KEYCLOAK_CLIENT_ID")

	router.Use(middleware.RateLimit())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.Sanitize())
	router.Use(middleware.Validate())

	// Instantiate the service and controller
	artworkRepo := repositories.NewGormArtworkRepository(initializers.DB) // GORM-based repository
	artworkService := services.NewArtworkService(artworkRepo)             // Service depends on the repository interface
	artworkController := controllers.NewArtworkController(artworkService) // Controller depends on the service

	galleryGroup := router.Group("/gallery")
	// galleryGroup.Use(middleware.Auth("Gallery", clientID))
	// galleryGroup.Use(middleware.Sanitize())
	// galleryGroup.Use(middleware.Validate())

	galleryGroup.POST("/artworks", artworkController.Create)
	galleryGroup.PUT("/artworks/:id", artworkController.Update)
	galleryGroup.DELETE("/artworks/:id", artworkController.Delete)

	regularGroup := router.Group("/")
	// regularGroup.Use(middleware.Auth("Regular", clientID))
	// regularGroup.Use(middleware.Sanitize())
	// regularGroup.Use(middleware.Validate())

	regularGroup.GET("/artworks", artworkController.Index)
	regularGroup.GET("/artworks/:id", artworkController.Find)

	router.GET("/auth/callback", middleware.Sanitize(), middleware.Validate(), controllers.AuthCallback)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
