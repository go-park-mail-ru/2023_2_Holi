package domain

type FavouritesRepository interface {
	Add(videoID, userID int) error
	Delete(videoID, userID int) error
	GetAll(userID int) ([]Video, error)
}

type FavouritesUsecase interface {
	Add(videoID, userID int) error
	Delete(videoID, userID int) error
	GetAll(userID int) ([]Video, error)
}
