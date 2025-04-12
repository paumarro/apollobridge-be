package integration_test

import (
	"log"
	"testing"
	"time"

	"github.com/paumarro/apollo-be/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Helper function to set up a fresh test database for each test
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Artwork{}); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}
	return db
}

// Helper function to clean up the test database
func teardownTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	sqlDB.Close()
}

// Test for creating an artwork
func TestArtworkCreate(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	date := time.Now()
	artwork := models.Artwork{
		Title:       "Starry Night",
		Artist:      "Vincent van Gogh",
		Date:        &date,
		Description: "A famous painting by van Gogh",
		Image:       "starry_night.jpg",
	}

	// Test Create
	if err := db.Create(&artwork).Error; err != nil {
		t.Fatalf("Failed to create artwork: %v", err)
	}

	// Verify Create
	var fetchedArtwork models.Artwork
	if err := db.First(&fetchedArtwork, artwork.ID).Error; err != nil {
		t.Fatalf("Failed to fetch created artwork: %v", err)
	}
	if fetchedArtwork.Title != artwork.Title || fetchedArtwork.Artist != artwork.Artist {
		t.Errorf("Expected artwork %+v, got %+v", artwork, fetchedArtwork)
	}
}

// Test for updating an artwork
func TestArtworkUpdate(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	// Create a new artwork
	date := time.Now()
	artwork := models.Artwork{
		Title:       "Starry Night",
		Artist:      "Vincent van Gogh",
		Date:        &date,
		Description: "A famous painting by van Gogh",
		Image:       "starry_night.jpg",
	}
	if err := db.Create(&artwork).Error; err != nil {
		t.Fatalf("Failed to create artwork: %v", err)
	}

	// Update the artwork
	newTitle := "The Starry Night"
	artwork.Title = newTitle
	if err := db.Save(&artwork).Error; err != nil {
		t.Fatalf("Failed to update artwork: %v", err)
	}

	// Verify Update
	var updatedArtwork models.Artwork
	if err := db.First(&updatedArtwork, artwork.ID).Error; err != nil {
		t.Fatalf("Failed to fetch updated artwork: %v", err)
	}
	if updatedArtwork.Title != newTitle {
		t.Errorf("Expected updated title %s, got %s", newTitle, updatedArtwork.Title)
	}
}

// Test for deleting an artwork
func TestArtworkDelete(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	// Create a new artwork
	date := time.Now()
	artwork := models.Artwork{
		Title:       "Starry Night",
		Artist:      "Vincent van Gogh",
		Date:        &date,
		Description: "A famous painting by van Gogh",
		Image:       "starry_night.jpg",
	}
	if err := db.Create(&artwork).Error; err != nil {
		t.Fatalf("Failed to create artwork: %v", err)
	}

	// Delete the artwork
	if err := db.Delete(&artwork).Error; err != nil {
		t.Fatalf("Failed to delete artwork: %v", err)
	}

	// Verify Delete
	var deletedArtwork models.Artwork
	if err := db.First(&deletedArtwork, artwork.ID).Error; err == nil {
		t.Fatalf("Expected artwork to be deleted, but it still exists")
	}
}

// Table-driven test for creating multiple artworks
func TestArtworkCreateMultiple(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	tests := []struct {
		name       string
		artwork    models.Artwork
		expectFail bool
	}{
		{
			name: "Valid Artwork",
			artwork: models.Artwork{
				Title:       "Starry Night",
				Artist:      "Vincent van Gogh",
				Description: "A famous painting by van Gogh",
				Image:       "starry_night.jpg",
			},
			expectFail: false,
		},
		{
			name: "Missing Title",
			artwork: models.Artwork{
				Artist:      "Vincent van Gogh",
				Description: "A famous painting by van Gogh",
				Image:       "starry_night.jpg",
			},
			expectFail: true,
		},
		{
			name: "Missing Artist",
			artwork: models.Artwork{
				Title:       "Starry Night",
				Description: "A famous painting by van Gogh",
				Image:       "starry_night.jpg",
			},
			expectFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Create(&tt.artwork).Error
			if (err != nil) != tt.expectFail {
				t.Fatalf("Unexpected result for %s: %v", tt.name, err)
			}
		})
	}
}
