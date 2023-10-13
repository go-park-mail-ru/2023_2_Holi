package genre_usecase

import "2023_2_Holi/domain"

type genreUsecase struct {
	genreRepo domain.GenreRepository
}

func NewGenreUsecase(gr domain.GenreRepository) domain.GenreUsecase {
	return &genreUsecase{
		genreRepo: gr,
	}
}

func (u *genreUsecase) GetGenres() ([]domain.Genre, error) {
	genres, err := u.genreRepo.GetGenres()
	if err != nil {
		return nil, err
	}

	return genres, nil
}
