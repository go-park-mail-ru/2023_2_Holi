package domain

type FavouritesRepository interface {
	Insert(videoID, userID int) error
	Delete(videoID, userID int) error
	SelectAll(userID int) ([]Video, error)
}

type FavouritesUsecase interface {
	Add(videoID, userID int) error
	Remove(videoID, userID int) error
	GetAll(userID int) ([]Video, error)
}
