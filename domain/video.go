package domain

import "github.com/jackc/pgx/v5/pgtype"

type Video struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	PreviewPath      string          `json:"previewPath"`
	MediaPath        string          `json:"mediaPath"`
	ReleaseYear      int             `json:"releaseYear"`
	Rating           float64         `json:"rating"`
	AgeRestriction   int             `json:"ageRestriction"`
	Duration         pgtype.Interval `json:"duration"`
	PreviewVideoPath string          `json:"previewVideoPath"`
	SeasonsCount     string          `json:"seasonsCount,omitempty"`
}
