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
	keycloakSecrete := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	log.Printf("PORT: %s", port)
	log.Printf("PG_DB_URL: %s", artDbURL)
	log.Printf("KEYCLOAK_URL: %s", keycloakURL)
	log.Printf("KEYCLOAK_SECRET: %s", keycloakSecrete)

	// Example of error handling for critical variables
	if port == "" || artDbURL == "" || keycloakURL == "" {
		log.Fatalf("Error: Critical environment variables are not set")
	}

	// Log the successful loading of environment variables
	log.Println("Environment variables loaded successfully!")

	log.Printf("PGUSER: %s", os.Getenv("PGUSER"))
	log.Printf("POSTGRES_PASSWORD: %s", os.Getenv("POSTGRES_PASSWORD"))
	log.Printf("RAILWAY_TCP_PROXY_DOMAIN: %s", os.Getenv("RAILWAY_TCP_PROXY_DOMAIN"))
	log.Printf("RAILWAY_TCP_PROXY_PORT: %s", os.Getenv("RAILWAY_TCP_PROXY_PORT"))
	log.Printf("PGDATABASE: %s", os.Getenv("PGDATABASE"))

}
