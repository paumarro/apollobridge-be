package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func ConnectToTestDB() {
	// Build the DSN from environment variables
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PG_TEST_USER"),
		os.Getenv("PG_TEST_PASSWORD"),
		os.Getenv("PG_TEST_HOST"),
		os.Getenv("PG_TEST_PORT"),
		os.Getenv("PG_TEST_DB"),
	)

	log.Println("Connecting to test database with DSN:", dsn)

	// Attempt to establish the connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Assign the connection to the global DB variable
	TestDB = db

	// Test the connection by running a simple query
	var result int
	if err := TestDB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		log.Fatalf("Error executing test query: %v", err)
	}

	log.Println("Database connection established successfully!")
}
