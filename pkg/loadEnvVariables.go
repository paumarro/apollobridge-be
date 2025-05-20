package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// Check if running in a Docker container
	if os.Getenv("RUNNING_IN_DOCKER") == "true" {
		log.Println("Running inside Docker container. Skipping LoadEnvVariables.")
		return
	}

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	artDbURL := os.Getenv("ART_DB_URL")

	// Example of error handling for critical variables
	if port == "" || artDbURL == "" || keycloakURL == "" {
		log.Fatalf("Error: Critical environment variables are not set")
	}

	// Log the successful loading of environment variables
	log.Println("Environment variables loaded successfully!")
}
