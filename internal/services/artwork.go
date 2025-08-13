package services

import (
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/paumarro/apollo-be/internal/repositories"
)

type ArtworkService struct {
	Repo repositories.ArtworkRepository
}

func NewArtworkService(repo repositories.ArtworkRepository) *ArtworkService {
	return &ArtworkService{Repo: repo}
}

func (s *ArtworkService) CreateArtwork(artwork *models.Artwork) error {
	return s.Repo.Create(artwork)
}

func (s *ArtworkService) GetAllArtworks() ([]models.Artwork, error) {
	return s.Repo.FindAll()
}

func (s *ArtworkService) GetArtworkByID(id string) (*models.Artwork, error) {
	return s.Repo.FindByID(id)
}

func (s *ArtworkService) UpdateArtwork(artwork *models.Artwork) error {
	return s.Repo.Update(artwork)
}

func (s *ArtworkService) DeleteArtwork(id string) error {
	return s.Repo.Delete(id)
}
