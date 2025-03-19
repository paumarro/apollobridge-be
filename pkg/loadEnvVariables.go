package env

import (
	"log"
	"os"
)

func LoadEnvVariables() {
	// Access and log environment variables
	port := os.Getenv("PORT")
	artDbURL := os.Getenv("ART_DB_URL")
	keycloakURL := os.Getenv("KEYCLOAK_URL")

	log.Printf("PORT: %s", port)
	log.Printf("ART_DB_URL: %s", artDbURL)
	log.Printf("KEYCLOAK_URL: %s", keycloakURL)

	// Example of error handling for critical variables
	if port == "" || artDbURL == "" || keycloakURL == "" {
		log.Fatalf("Error: Critical environment variables are not set")
	}

	// Log the successful loading of environment variables
	log.Println("Environment variables loaded successfully!")
}
