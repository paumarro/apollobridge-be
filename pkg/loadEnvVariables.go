package env

import (
	"log"
	"os"
)

func LoadEnvVariables() {
	// Access and log environment variables
	port := os.Getenv("PORT")
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	artDbURL := os.Getenv("PG_DB_URL")

	// Example of error handling for critical variables
	if port == "" || artDbURL == "" || keycloakURL == "" {
		log.Fatalf("Error: Critical environment variables are not set")
	}

	// Log the successful loading of environment variables
	log.Println("Environment variables loaded successfully!")
}
