// internal/models/artwork.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Artwork struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Artist      string         `json:"artist"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
