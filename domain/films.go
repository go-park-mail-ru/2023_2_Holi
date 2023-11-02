package domain

import "github.com/jackc/pgx/v5/pgtype"

type Film struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"-"`
	PreviewPath      string          `json:"previewPath"`
	MediaPath        string          `json:"-"`
	ReleaseYear      int             `json:"-"`
	Rating           float64         `json:"rating"`
	AgeRestriction   int             `json:"-"`
	Duration         pgtype.Interval `json:"-"`
	PreviewVideoPath string          `json:"previewVideoPath"`
}

type Cast struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

type FilmsRepository interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (Film, error)
	GetFilmCast(filmId int) ([]Cast, error)
	GetCastPage(id int) ([]Film, error)
	GetCastName(id int) (Cast, error)
}

type FilmsUsecase interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (Film, []Cast, error)
	GetCastPage(id int) ([]Film, Cast, error)
}
