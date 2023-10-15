package domain

type Movie struct {
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

type MoviesRepository interface {
	GetMoviesByGenre(genre string) ([]Movie, error)
}

type MoviesUsecase interface {
	GetMoviesByGenre(genre string) ([]Movie, error)
}
