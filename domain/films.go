package domain

import "github.com/jackc/pgx/v5/pgtype"

type Film struct {
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
}

type Episode struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PreviewPath string `json:"previewPath"`
	MediaPath   string `json:"mediaPath"`
}

type Cast struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FilmsRepository interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (Film, error)
	GetFilmCast(filmId int) ([]Cast, error)
	GetCastPage(id int) ([]Film, error)
	GetCastName(id int) (Cast, error)
	GetTopRate() (Film, error)
}

type FilmsUsecase interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (Film, []Cast, error)
	GetCastPage(id int) ([]Film, Cast, error)
	GetTopRate() (Film, error)
}

type SeriesRepository interface {
	GetSeriesByGenre(genre string) ([]Film, error)
	GetSeriesData(id int) (Film, []Cast, error)
}

type SeriesUsecase interface {
	GetSeriesByGenre(genre string) ([]Film, error)
	GetSeriesData(id int) (Film, []Cast, error)
}
