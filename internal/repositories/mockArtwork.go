package repositories

import (
	"github.com/paumarro/apollo-be/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockArtworkRepository struct {
	mock.Mock
}

func (m *MockArtworkRepository) Create(artwork *models.Artwork) error {
	args := m.Called(artwork)
	return args.Error(0)
}

func (m *MockArtworkRepository) FindAll() ([]models.Artwork, error) {
	args := m.Called()
	var res []models.Artwork
	if v := args.Get(0); v != nil {
		res = v.([]models.Artwork)
	}
	return res, args.Error(1)
}

func (m *MockArtworkRepository) FindByID(id string) (*models.Artwork, error) {
	args := m.Called(id)
	var res *models.Artwork
	if v := args.Get(0); v != nil {
		res = v.(*models.Artwork)
	}
	return res, args.Error(1)
}

func (m *MockArtworkRepository) Update(artwork *models.Artwork) error {
	args := m.Called(artwork)
	return args.Error(0)
}

func (m *MockArtworkRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
