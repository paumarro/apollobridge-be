package dto

type ArtworkRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`  // Required, min 3, max 255 characters
	Artist      string `json:"artist" validate:"required,min=3,max=255"` // Required, min 3, max 255 characters
	Description string `json:"description" validate:"required,max=1000"` // Required, max 1000 characters
	Image       string `json:"image" validate:"required,url"`            // Required, must be a valid URL
}
