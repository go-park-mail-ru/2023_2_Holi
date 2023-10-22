package domain

type Artist struct {
	ID      int
	Name    string
	Surname string
}
type ArtistUsecase interface {
	GetArtistPage(name string) ([]Film, error)
}

type ArtistRepository interface {
	GetArtistPage(name string) ([]Film, error)
}
