package domain

type SearchData struct {
	Cast  []Cast
	Films []Film
}

type SearchUsecase interface {
	GetSearchData(searchStr string) (SearchData, error)
}

type SearchRepository interface {
	GetSuitableFilms(searchStr string) ([]Film, error)
	GetSuitableCast(searchStr string) ([]Cast, error)
}
