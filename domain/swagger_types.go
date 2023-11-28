// Package domain is used for swagger auto doc
package domain

type UserRequest struct {
	ID        int    `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	ImagePath string `json:"imagePath"`
}

type VideoResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	PreviewPath      string  `json:"previewPath"`
	MediaPath        string  `json:"mediaPath"`
	ReleaseYear      int     `json:"releaseYear"`
	Rating           float64 `json:"rating"`
	AgeRestriction   int     `json:"ageRestriction"`
	Duration         string  `json:"duration"`
	PreviewVideoPath string  `json:"previewVideoPath"`
}
