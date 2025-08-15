package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
)

// Predefined errors for the service layer
var (
	ErrNotFound = errors.New("record not found")
)

// ArtworkService provides business logic for managing artworks.
type ArtworkService struct {
	Repo repositories.ArtworkRepository
}

// NewArtworkService creates a new instance of ArtworkService.
func NewArtworkService(repo repositories.ArtworkRepository) *ArtworkService {
	return &ArtworkService{Repo: repo}
}

// CreateArtwork creates a new artwork in the repository.
func (s *ArtworkService) CreateArtwork(artwork *models.Artwork) error {
	log.Printf("Creating artwork: %+v", artwork)

	// Call the repository to create the artwork
	if err := s.Repo.Create(artwork); err != nil {
		log.Printf("Error in CreateArtwork: %v", err)
		return fmt.Errorf("failed to create artwork: %w", err) // Wrap the error with context
	}

	return nil
}

// GetAllArtworks retrieves all artworks from the repository.
func (s *ArtworkService) GetAllArtworks() ([]models.Artwork, error) {
	log.Println("Fetching all artworks")

	// Call the repository to fetch all artworks
	artworks, err := s.Repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artworks: %w", err) // Wrap the error with context
	}

	return artworks, nil
}

// GetArtworkByID retrieves a single artwork by ID.
func (s *ArtworkService) GetArtworkByID(id string) (*models.Artwork, error) {
	log.Printf("Fetching artwork with ID: %s", id)

	// Call the repository to fetch the artwork by ID
	artwork, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound // Return a predefined "not found" error
		}
		return nil, fmt.Errorf("failed to fetch artwork with ID %s: %w", id, err) // Wrap the error with context
	}

	return artwork, nil
}

// UpdateArtwork updates an artwork in the repository.
func (s *ArtworkService) UpdateArtwork(artwork *models.Artwork) error {
	log.Printf("Updating artwork with ID: %d", artwork.ID)

	// Call the repository to update the artwork
	if err := s.Repo.Update(artwork); err != nil {
		return fmt.Errorf("failed to update artwork with ID %d: %w", artwork.ID, err) // Wrap the error with context
	}

	return nil
}

// DeleteArtwork deletes an artwork by ID.
func (s *ArtworkService) DeleteArtwork(id string) error {
	log.Printf("Deleting artwork with ID: %s", id)

	// Call the repository to delete the artwork
	if err := s.Repo.Delete(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound // Return a predefined "not found" error
		}
		return fmt.Errorf("failed to delete artwork with ID %s: %w", id, err) // Wrap the error with context
	}

	return nil
}
