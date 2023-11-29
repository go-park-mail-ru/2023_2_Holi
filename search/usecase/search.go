package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type searchUseCase struct {
	searchRepo domain.SearchRepository
}

func NewSearchUsecase(sr domain.SearchRepository) domain.SearchUsecase {
	return &searchUseCase{searchRepo: sr}
}

func (u *searchUseCase) GetSearchData(searchStr string) (domain.SearchData, error) {
	films, err := u.searchRepo.GetSuitableFilms(searchStr)
	if err != nil {
		logs.LogError(logs.Logger, "search_usecase", "GetSearchData", err, err.Error())
	}
	logs.Logger.Debug("Usecase GetSearchData films:", films)

	cast, err := u.searchRepo.GetSuitableCast(searchStr)
	if err != nil {
		logs.LogError(logs.Logger, "search_usecase", "GetSearchData", err, err.Error())
	}
	logs.Logger.Debug("Usecase GetSearchData cast:", cast)

	if len(films) == 0 && len(cast) == 0 {
		return domain.SearchData{}, domain.ErrNotFound
	}

	data := domain.SearchData{
		Cast:  cast,
		Films: films,
	}

	return data, nil
}
