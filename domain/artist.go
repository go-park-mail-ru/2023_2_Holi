package domain

type Artist struct {
	ID      int
	Name    string
	Surname string
}

type ArtistUsecase interface {
	GetArtistPage(name, surname string) ([]Film, error)
}

type ArtistRepository interface {
	GetArtistPage(name, surname string) ([]Film, error)
}
