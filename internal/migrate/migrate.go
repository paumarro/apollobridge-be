package main

import (
	"log"

	initializers "github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/models"
	env "github.com/paumarro/apollo-be/pkg"
)

func init() {
	env.LoadEnvVariables(".../.../.env")
	initializers.ConnectToDB()
}

func main() {
	if err := initializers.DB.AutoMigrate(&models.Artwork{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
