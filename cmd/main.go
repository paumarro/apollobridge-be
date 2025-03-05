package main

import (
	"log"

	"github.com/paumarro/apollo-be/internal/config"
	"github.com/paumarro/apollo-be/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create and start server
	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
