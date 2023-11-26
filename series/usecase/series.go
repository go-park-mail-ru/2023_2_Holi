package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type seriesUsecase struct {
	seriesRepo domain.SeriesRepository
}

func NewFilmsUsecase(sr domain.SeriesRepository) domain.SeriesUsecase {
	return &seriesUsecase{
		seriesRepo: sr,
	}
}

func (u *seriesUsecase) GetSeriesByGenre(genre string) ([]domain.Video, error) {
	films, err := u.seriesRepo.GetSeriesByGenre(genre)
	if err != nil {
		logs.LogError(logs.Logger, "films_usecase", "GetFilmsByGenre", err, err.Error())
		return nil, err
	}
	logs.Logger.Debug("Usecase GetFilmsByGenre:", films)
	return films, nil
}

func (u *seriesUsecase) GetSeriesData(id int) (domain.Video, []domain.Cast, []domain.Episode, error) {
	film, err := u.seriesRepo.GetSeriesData(id)
	if err != nil {
		logs.LogError(logs.Logger, "films_usecase", "GetFilmData", err, err.Error())
		return domain.Video{}, nil, nil, err
	}
	artists, err := u.seriesRepo.GetSeriesCast(id)
	if err != nil {
		return domain.Video{}, []domain.Cast{}, nil, err
	}
	episodes, err := u.seriesRepo.GetSeriesEpisodes(id)
	if err != nil {
		return domain.Video{}, nil, []domain.Episode{}, err
	}

	return film, artists, episodes, nil
}

func (u *seriesUsecase) GetCastPageSeries(id int) ([]domain.Video, domain.Cast, error) {
	series, err := u.seriesRepo.GetCastPageSeries(id)
	if err != nil {
		return nil, domain.Cast{}, err
	}
	logs.Logger.Debug("Usecase GetCastPageSeries:", series)
	artist, err := u.seriesRepo.GetCastNameSeries(id)
	if err != nil {
		return nil, domain.Cast{}, err
	}
	logs.Logger.Debug("Usecase GetCastPageSeries:", artist)
	return series, artist, nil
}

//func (u *seriesUsecase) GetTopRate() (domain.Video, error) {
//	genres, err := u.seriesRepo.GetTopRate()
//	if err != nil {
//		return domain.Video{}, err
//	}
//	//logger.Debug("Usecase GetGenres:", genres)
//
//	return genres, nil
//}
