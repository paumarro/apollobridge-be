package models

import (
	"time"

	"gorm.io/gorm"
)

type Artwork struct {
	gorm.Model
	Title       string
	Artist      string
	Date        *time.Time
	Description string
}
