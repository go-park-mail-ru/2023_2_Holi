package domain

type Film struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"-"`
	PreviewPath    string  `json:"previewPath"`
	MediaPath      string  `json:"-"`
	Country        string  `json:"-"`
	ReleaseYear    int     `json:"-"`
	Rating         float64 `json:"rating"`
	RatesCount     int     `json:"-"`
	AgeRestriction int     `json:"-"`
	Duration       int     `json:"-"`
}

type FilmsRepository interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (*Film, error)
	GetFilmArtists(filmId int) ([]Artist, error)
}

type FilmsUsecase interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (*Film, []Artist, error)
}
