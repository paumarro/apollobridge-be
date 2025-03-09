package main

import (
	initializers "github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/models"
	env "github.com/paumarro/apollo-be/pkg"
)

func init() {
	env.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Artwork{})
}
