package domain

type Genre struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

type GenreRepository interface {
	GetGenres() ([]Genre, error)
	//GetGenresSeries() ([]Genre, error)
}

type GenreUsecase interface {
	GetGenres() ([]Genre, error)
	//GetGenresSeries() ([]Genre, error)
}
