package genre_postgres

import (
	"context"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

var logger = logs.LoggerInit()

type genrePostgresRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func GenrePostgresqlRepository(pool *pgxpool.Pool, ctx context.Context) domain.GenreRepository {
	return &genrePostgresRepo{
		db:  pool,
		ctx: ctx,
	}
}

func (r *genrePostgresRepo) GetGenres() ([]domain.Genre, error) {
	sql, args, err := domain.Psql.Select("name").
		From("genre").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(r.ctx, sql, args...)
	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "genre_postgres", "GetGenres", err, err.Error())
		return nil, domain.ErrNotFound
	}
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

	return genres, nil
}
