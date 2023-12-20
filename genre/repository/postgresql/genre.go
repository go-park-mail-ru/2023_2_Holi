package postgres

import (
	"context"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const getGenresFilmsQuery = `
	SELECT id, name
	FROM genre
`

const getGenresSeriesQuery = `
	SELECT DISTINCT g.id, g.name
	FROM genre AS g
	JOIN video_genre AS vg ON g.id = vg.genre_id
	JOIN video AS v ON vg.video_id = v.id
	WHERE v.seasons_count > 0;
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
	rows, err := r.db.Query(r.ctx, getGenresFilmsQuery)
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
			&genre.ID,
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

func (r *genrePostgresRepo) GetGenresSeries() ([]domain.Genre, error) {
	rows, err := r.db.Query(r.ctx, getGenresSeriesQuery)
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
			&genre.ID,
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
