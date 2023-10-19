package domain

type ArtistUsecase interface {
	GetArtistPage(name, surname string) ([]Film, error)
}

type ArtistRepository interface {
	GetArtistPage(name, surname string) ([]Film, error)
}
