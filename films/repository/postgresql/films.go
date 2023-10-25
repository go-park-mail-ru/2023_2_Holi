package films_postgres

import (
	"context"
	"github.com/jackc/pgx/v5"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type filmsPostgresqlRepository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewFilmsPostgresqlRepository(pool *pgxpool.Pool, ctx context.Context) domain.FilmsRepository {
	return &filmsPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *filmsPostgresqlRepository) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	sqlString, args, err := domain.Psql.Select("video.id", "video.name", "video.preview_path", "video.rating").
		From("video").
		Join("video_genre AS vg ON video_id = vg.video_id").
		Join("genre AS g ON vg.genre_id = g.id").
		Where("g.name = ?", genre).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(r.ctx, sqlString, args...)
	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return nil, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmsByGenre query result:", rows)

	var films []domain.Film
	for rows.Next() {
		var film domain.Film
		err = rows.Scan(
			&film.ID,
			&film.Name,
			&film.PreviewPath,
			&film.Rating,
		)

		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
			return nil, err
		}

		films = append(films, film)
	}

	return films, nil
}

func (r *filmsPostgresqlRepository) GetFilmData(id int) (*domain.Film, error) {
	sql, args, err := domain.Psql.Select("e.name", "e.description", "e.duration",
		"e.preview_path", "e.media_path", "release_year", "rating", "age_restriction").
		From("video").
		Join("episode AS e ON video.id = video_id").
		Where("video.id = ?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	row, err := r.db.Query(r.ctx, sql, args...)
	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return nil, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return nil, err
	}
	defer row.Close()
	logs.Logger.Debug("GetFilmData query result:", row)

	film := new(domain.Film)
	err = row.Scan(
		&film.Name,
		&film.Description,
		&film.Duration,
		&film.PreviewPath,
		&film.MediaPath,
		&film.ReleaseYear,
		&film.Rating,
		&film.AgeRestriction,
	)

	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return nil, err
	}

	return film, nil
}

func (r *filmsPostgresqlRepository) GetFilmArtists(FilmId int) ([]domain.Artist, error) {
	sql, args, err := domain.Psql.Select("name").
		From("cast").
		Join("video_film AS vf ON id = cast_id").
		Where("af.video_id = ?", FilmId).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(r.ctx, sql, args...)
	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmArtists", err, err.Error())
		return nil, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmArtists", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmArtists query result:", rows)

	var artists []domain.Artist
	for rows.Next() {
		var artist domain.Artist
		err = rows.Scan(
			&artist.Name,
		)
		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmArtists", err, err.Error())
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, nil
}
