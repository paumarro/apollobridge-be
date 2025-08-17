package services

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
)

// Predefined errors for the service layer
var (
	ErrNotFound   = errors.New("record not found")
	ErrBadRequest = errors.New("bad request")
	ErrDuplicate  = errors.New("Artwork already exists")
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

	// Business rule: prevent duplicates by title + artist
	existing, err := s.Repo.FindAll()
	if err != nil {
		log.Printf("Error checking duplicates: %v", err)
		return fmt.Errorf("failed to create artwork: %w", err)
	}
	for _, a := range existing {
		if a.Title == artwork.Title && a.Artist == artwork.Artist {
			return ErrDuplicate
		}
	}

	// Persist
	if err := s.Repo.Create(artwork); err != nil {
		log.Printf("Error in CreateArtwork: %v", err)
		return fmt.Errorf("failed to create artwork: %w", err)
	}
	return nil
}

// GetAllArtworks retrieves all artworks from the repository.
func (s *ArtworkService) GetAllArtworks() ([]models.Artwork, error) {
	log.Println("Fetching all artworks")
	artworks, err := s.Repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artworks: %w", err)
	}
	return artworks, nil
}

// GetArtworkByID retrieves a single artwork by ID.
func (s *ArtworkService) GetArtworkByID(id string) (*models.Artwork, error) {
	log.Printf("Fetching artwork with ID: %s", id)
	artwork, err := s.Repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch artwork with ID %s: %w", id, err)
	}
	return artwork, nil
}

// UpdateArtwork updates an artwork in the repository.
func (s *ArtworkService) UpdateArtwork(artwork *models.Artwork) error {
	log.Printf("Updating artwork with ID: %d", artwork.ID)
	if err := s.Repo.Update(artwork); err != nil {
		return fmt.Errorf("failed to update artwork with ID %d: %w", artwork.ID, err)
	}
	return nil
}

// DeleteArtwork deletes an artwork by ID.
func (s *ArtworkService) DeleteArtwork(id string) error {
	log.Printf("Deleting artwork with ID: %s", id)
	if err := s.Repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete artwork with ID %s: %w", id, err)
	}
	return nil
}
