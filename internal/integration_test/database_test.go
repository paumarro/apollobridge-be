package integration_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/paumarro/apollo-be/internal/initializers"
	"github.com/paumarro/apollo-be/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Set up a test database
	initializers.DB = setupTestDB()

	// Run migrations
	err := initializers.DB.AutoMigrate(&models.Artwork{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Clean up
	teardownTestDB()

	os.Exit(code)
}

func setupTestDB() *gorm.DB {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func teardownTestDB() {
	sqlDB, err := initializers.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	sqlDB.Close()
}

func TestArtworkCRUD(t *testing.T) {
	// Create a new artwork
	date := time.Now()
	artwork := models.Artwork{
		Title:       "Starry Night",
		Artist:      "Vincent van Gogh",
		Date:        &date,
		Description: "A famous painting by van Gogh",
		Image:       "starry_night.jpg",
	}

	// Test Create
	if err := initializers.DB.Create(&artwork).Error; err != nil {
		t.Fatalf("Failed to create artwork: %v", err)
	}

	// Verify Create
	var fetchedArtwork models.Artwork
	if err := initializers.DB.First(&fetchedArtwork, artwork.ID).Error; err != nil {
		t.Fatalf("Failed to fetch created artwork: %v", err)
	}
	if fetchedArtwork.Title != artwork.Title {
		t.Errorf("Expected title %s, got %s", artwork.Title, fetchedArtwork.Title)
	}

	// Test Update
	newTitle := "The Starry Night"
	fetchedArtwork.Title = newTitle
	if err := initializers.DB.Save(&fetchedArtwork).Error; err != nil {
		t.Fatalf("Failed to update artwork: %v", err)
	}

	// Verify Update
	var updatedArtwork models.Artwork
	if err := initializers.DB.First(&updatedArtwork, artwork.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated artwork: %v", err)
	}
	if updatedArtwork.Title != newTitle {
		t.Errorf("Expected updated title %s, got %s", newTitle, updatedArtwork.Title)
	}

	// Test Delete
	if err := initializers.DB.Delete(&updatedArtwork).Error; err != nil {
		t.Fatalf("Failed to delete artwork: %v", err)
	}

	// Verify Delete
	var deletedArtwork models.Artwork
	if err := initializers.DB.First(&deletedArtwork, artwork.ID).Error; err == nil {
		t.Fatalf("Expected artwork to be deleted, but it still exists")
	}
}
