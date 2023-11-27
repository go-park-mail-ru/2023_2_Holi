package postgres

import (
	"context"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const getGenresQuery = `
	SELECT name
	FROM genre
`

type genrePostgresRepo struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func GenrePostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.GenreRepository {
	return &genrePostgresRepo{
		db:  pool,
		ctx: ctx,
	}
}

func (r *genrePostgresRepo) GetGenres() ([]domain.Genre, error) {
	rows, err := r.db.Query(r.ctx, getGenresQuery)
	if err != nil {
		logs.LogError(logs.Logger, "genre_postgres", "GetGenres", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetGenres query result:", rows)

	var genres []domain.Genre
	for rows.Next() {
		var genre domain.Genre
		err = rows.Scan(
			&genre.Name,
		)

		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}
	logs.Logger.Info("lenth", len(genres))
	if len(genres) == 0 {
		return nil, domain.ErrNotFound
	}

	return genres, nil
}

//func (r *genrePostgresRepo) GetGenresSeries() ([]domain.Genre, error) {
//	rows, err := r.db.Query(r.ctx, getGenresQuery)
//	if err != nil {
//		logs.LogError(logs.Logger, "genre_postgres", "GetGenres", err, err.Error())
//		return nil, err
//	}
//	defer rows.Close()
//	logs.Logger.Debug("GetGenres query result:", rows)
//
//	var genres []domain.Genre
//	for rows.Next() {
//		var genre domain.Genre
//		err = rows.Scan(
//			&genre.Name,
//		)
//
//		if err != nil {
//			return nil, err
//		}
//
//		genres = append(genres, genre)
//	}
//	logs.Logger.Info("lenth", len(genres))
//	if len(genres) == 0 {
//		return nil, domain.ErrNotFound
//	}
//
//	return genres, nil
//}
