package models

import (
	"time"

	"gorm.io/gorm"
)

type Artwork struct {
	gorm.Model
	ID          uint
	Title       string
	Artist      string
	Date        *time.Time
	Description string
	Image       string
}
