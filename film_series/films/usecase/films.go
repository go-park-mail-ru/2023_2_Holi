package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type filmsUsecase struct {
	filmRepo domain.FilmsRepository
	recRepo  domain.RecomRepository
}

func NewFilmsUsecase(fr domain.FilmsRepository, fc domain.RecomRepository) domain.FilmsUsecase {
	return &filmsUsecase{
		filmRepo: fr,
		recRepo:  fc,
	}
}

func (u *filmsUsecase) GetFilmsByGenre(genre int) ([]domain.Video, error) {
	films, err := u.filmRepo.GetFilmsByGenre(genre)
	if err != nil {
		logs.LogError(logs.Logger, "films_usecase", "GetFilmsByGenre", err, err.Error())
		return nil, err
	}
	logs.Logger.Debug("Usecase GetFilmsByGenre:", films)
	return films, nil
}

func (u *filmsUsecase) GetCastPage(id int) ([]domain.Video, domain.Cast, error) {
	films, err := u.filmRepo.GetCastPage(id)
	if err != nil {
		return nil, domain.Cast{}, err
	}
	logs.Logger.Debug("Usecase GetCastPage:", films)
	artist, err := u.filmRepo.GetCastName(id)
	if err != nil {
		return nil, domain.Cast{}, err
	}
	logs.Logger.Debug("Usecase GetCastPage:", artist)
	return films, artist, nil
}

func (u *filmsUsecase) GetFilmData(id int) (domain.Video, []domain.Cast, error) {
	film, err := u.filmRepo.GetFilmData(id)
	if err != nil {
		logs.LogError(logs.Logger, "films_usecase", "GetFilmData", err, err.Error())
		return domain.Video{}, nil, err
	}
	artists, err := u.filmRepo.GetFilmCast(id)
	if err != nil {
		return domain.Video{}, nil, err
	}

	return film, artists, nil
}

func (u *filmsUsecase) GetTopRate() (domain.Video, error) {
	topRate, err := u.filmRepo.GetTopRate()
	if err != nil {
		return domain.Video{}, err
	}
	logs.Logger.Debug("films_usecase GetTopRate:", topRate)

	return topRate, nil
}

func (u *filmsUsecase) GetRecommendations(userID int) ([]domain.Video, error) {
	recommendations, err := u.recRepo.GetRecommendations(userID)
	if err != nil {
		return nil, err
	}
	logs.Logger.Debug("films_usecase GetRecommendations:", recommendations)

	var films []domain.Video
	for _, id := range recommendations {
		film, err := u.filmRepo.GetFilmData(id)
		if err != nil {
			logs.LogError(logs.Logger, "films_usecase", "GetFilmData", err, err.Error())
			return nil, err
		}

		films = append(films, film)
	}

	return films, nil
}
