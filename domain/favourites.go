package domain

type FavouritesRepository interface {
	InsertIntoFavourites(videoID, userID int) error
	DeleteFromFavourites(videoID, userID int) error
	SelectAllFavourites(userID int) ([]Video, error)
	Exists(videoID, userID int) (bool, error)
}

type FavouritesUsecase interface {
	AddToFavourites(videoID, userID int) error
	RemoveFromFavourites(videoID, userID int) error
	GetAllFavourites(userID int) ([]Video, error)
	Favourite(videoID, userID int) (bool, error)
}
