package models

import (
	"time"

	"gorm.io/gorm"
)

type Artwork struct {
	gorm.Model
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Artist      string     `json:"artist"`
	Date        *time.Time `json:"date,omitempty"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
}
