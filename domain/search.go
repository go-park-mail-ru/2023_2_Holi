package domain

type SearchData struct {
	Cast  []Cast
	Films []Video
}

type SearchUsecase interface {
	GetSearchData(searchStr string) (SearchData, error)
}

type SearchRepository interface {
	GetSuitableFilms(searchStr string) ([]Video, error)
	GetSuitableCast(searchStr string) ([]Cast, error)
}
