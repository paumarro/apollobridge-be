package services

import (
	"github.com/paumarro/apollo-be/internal/models"
	"gorm.io/gorm"
)

type ArtworkService struct {
	DB *gorm.DB
}

func NewArtworkService(db *gorm.DB) *ArtworkService {
	return &ArtworkService{DB: db}
}

func (s *ArtworkService) CreateArtwork(artwork *models.Artwork) error {
	return s.DB.Create(artwork).Error
}

func (s *ArtworkService) GetAllArtworks() ([]models.Artwork, error) {
	var artworks []models.Artwork
	err := s.DB.Find(&artworks).Error
	return artworks, err
}

func (s *ArtworkService) GetArtworkByID(id string) (*models.Artwork, error) {
	var artwork models.Artwork
	err := s.DB.First(&artwork, id).Error
	return &artwork, err
}

func (s *ArtworkService) UpdateArtwork(artwork *models.Artwork) error {
	return s.DB.Save(artwork).Error
}

func (s *ArtworkService) DeleteArtwork(id string) error {
	return s.DB.Delete(&models.Artwork{}, id).Error
}
