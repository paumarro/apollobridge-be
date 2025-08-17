package repositories

import (
	"github.com/paumarro/apollo-be/internal/models"
	"gorm.io/gorm"
)

// ArtworkRepository defines the database operations for artworks
type ArtworkRepository interface {
	Create(artwork *models.Artwork) error
	FindAll() ([]models.Artwork, error)
	FindByID(id string) (*models.Artwork, error)
	Update(artwork *models.Artwork) error
	Delete(id string) error
}

// GormArtworkRepository is the GORM-based implementation of ArtworkRepository
type GormArtworkRepository struct {
	DB *gorm.DB
}

func NewGormArtworkRepository(db *gorm.DB) *GormArtworkRepository {
	return &GormArtworkRepository{DB: db}
}

func (r *GormArtworkRepository) Create(artwork *models.Artwork) error {
	return r.DB.Create(artwork).Error
}

func (r *GormArtworkRepository) FindAll() ([]models.Artwork, error) {
	var artworks []models.Artwork
	err := r.DB.Find(&artworks).Error
	return artworks, err
}

func (r *GormArtworkRepository) FindByID(id string) (*models.Artwork, error) {
	var artwork models.Artwork
	if err := r.DB.First(&artwork, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &artwork, nil
}

func (r *GormArtworkRepository) Update(artwork *models.Artwork) error {
	return r.DB.Save(artwork).Error
}

func (r *GormArtworkRepository) Delete(id string) error {
	res := r.DB.Delete(&models.Artwork{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
