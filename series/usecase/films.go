package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type seriesUsecase struct {
	seriesRepo domain.SeriesRepository
}

func NewFilmsUsecase(sr domain.SeriesRepository) domain.SeriesRepository {
	return &seriesUsecase{
		seriesRepo: sr,
	}
}

func (u *seriesUsecase) GetSeriesByGenre(genre string) ([]domain.Film, error) {
	films, err := u.seriesRepo.GetSeriesByGenre(genre)
	if err != nil {
		logs.LogError(logs.Logger, "films_usecase", "GetFilmsByGenre", err, err.Error())
		return nil, err
	}
	logs.Logger.Debug("Usecase GetFilmsByGenre:", films)
	return films, nil
}

// func (u *filmsUsecase) GetCastPage(id int) ([]domain.Film, domain.Cast, error) {
// 	films, err := u.filmRepo.GetCastPage(id)
// 	if err != nil {
// 		return nil, domain.Cast{}, err
// 	}
// 	logs.Logger.Debug("Usecase GetCastPage:", films)
// 	artist, err := u.filmRepo.GetCastName(id)
// 	if err != nil {
// 		return nil, domain.Cast{}, err
// 	}
// 	logs.Logger.Debug("Usecase GetCastPage:", artist)
// 	return films, artist, nil
// }

func (u *seriesUsecase) GetSeriesData(id int) (domain.Film, []domain.Cast, error) {
	film, err := u.seriesRepo.GetSeriesData(id)
	if err != nil {
		logs.LogError(logs.Logger, "films_usecase", "GetFilmData", err, err.Error())
		return domain.Film{}, nil, err
	}
	artists, err := u.seriesRepo.GetSeriesCast(id)
	if err != nil {
		return domain.Film{}, nil, err
	}

	return film, artists, nil
}

// func (u *filmsUsecase) GetTopRate() (domain.Film, error) {
// 	genres, err := u.filmRepo.GetTopRate()
// 	if err != nil {
// 		return domain.Film{}, err
// 	}
// 	//logger.Debug("Usecase GetGenres:", genres)

// 	return genres, nil
// }
