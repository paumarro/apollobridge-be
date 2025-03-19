package env

import (
	"log"
	"os"
)

func LoadEnvVariables() {
	// Access environment variables
	port := os.Getenv("PORT")
	artDbURL := os.Getenv("PG_DB_URL")
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	// keycloakRealm := os.Getenv("KEYCLOAK_REALM")
	// keycloakClientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	// keycloakClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	// redirectURL := os.Getenv("REDIRECT_URL")

	// Example of error handling for critical variables
	if port == "" || artDbURL == "" || keycloakURL == "" {
		log.Fatalf("Error: Critical environment variables are not set")
	}

	// Log the successful loading of environment variables
	log.Println("Environment variables loaded successfully!")
}
