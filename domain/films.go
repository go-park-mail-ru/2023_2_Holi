package domain

type Film struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"-"`
	PreviewPath    string  `json:"previewPath"`
	MediaPath      string  `json:"-"`
	ReleaseYear    int     `json:"-"`
	Rating         float64 `json:"rating"`
	AgeRestriction int     `json:"-"`
	Duration       int     `json:"-"`
}

type Cast struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FilmsRepository interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (*Film, error)
	GetFilmCast(filmId int) ([]Cast, error)
	GetCastPage(id int) ([]Film, error)
	GetCastName(id int) ([]Cast, error)
}

type FilmsUsecase interface {
	GetFilmsByGenre(genre string) ([]Film, error)
	GetFilmData(id int) (*Film, []Cast, error)
	GetCastPage(id int) ([]Film, []Cast, error)
}
