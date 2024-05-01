package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const getFilmsByGenreQuery = `
	SELECT DISTINCT v.id, v.name, v.preview_path, v.rating, v.preview_video_path, v.seasons_count
	FROM video AS v
		JOIN video_cast AS vc ON v.id = vc.video_id
		JOIN "cast" AS c ON vc.cast_id = c.id
		JOIN episode AS e ON e.video_id = v.id
		JOIN video_genre AS vg ON v.id = vg.video_id
		JOIN genre AS g ON vg.genre_id = g.id
	WHERE g.id = $1 AND v.seasons_count = 0;
`

const getFilmDataQuery = `
    SELECT e.name, e.description, e.duration,
        e.preview_path, e.media_path, preview_video_path, release_year, rating, age_restriction
    FROM video
        JOIN episode AS e ON video.id = e.video_id
    WHERE video.id = $1 AND video.seasons_count = 0;
`

const getFilmCastQuery = `
    SELECT id, name, imgpath
    FROM "cast"
        JOIN video_cast AS vc ON id = cast_id
    WHERE vc.video_id = $1;
`

const getCastPageQuery = `
    SELECT video.id, e.name, e.preview_path, video.rating, video.preview_video_path
    FROM video
    JOIN episode AS e ON video.id = e.video_id
    WHERE video.id IN (
        SELECT vc.video_id
        FROM video_cast AS vc
        WHERE vc.cast_id = $1 AND video.seasons_count = 0
    );
`

const getCastNameQuery = `
    SELECT "cast".name, "cast".birthday, "cast".place, "cast".carier, "cast".imgpath
    FROM "cast" 
    WHERE "cast".id = $1;
`

const getTopRateQuery = `
    SELECT video.id, video.name, video.description, video.preview_video_path, e.media_path
    FROM video
    JOIN episode AS e ON video.id = e.video_id
    WHERE video.seasons_count = 0
    ORDER BY rating DESC
    LIMIT 1;
`

type filmsPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewFilmsPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.FilmsRepository {
	return &filmsPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *filmsPostgresqlRepository) GetFilmsByGenre(genre int) ([]domain.Video, error) {
	rows, err := r.db.Query(r.ctx, getFilmsByGenreQuery, genre)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmsByGenre query result:", rows)

	var films []domain.Video

	for rows.Next() {
		var film domain.Video
		err = rows.Scan(
			&film.ID,
			&film.Name,
			&film.PreviewPath,
			&film.Rating,
			&film.PreviewVideoPath,
			&film.SeasonsCount,
		)

		if err != nil {
			return nil, err
		}

		films = append(films, film)
	}

	logs.Logger.Info("lenth", len(films))
	if len(films) == 0 {
		return nil, domain.ErrNotFound
	}

	return films, nil
}

func (r *filmsPostgresqlRepository) GetFilmData(id int) (domain.Video, error) {
	row := r.db.QueryRow(r.ctx, getFilmDataQuery, id)

	logs.Logger.Debug("GetFilmData query result:", row)

	var film domain.Video
	err := row.Scan(
		&film.Name,
		&film.Description,
		&film.Duration,
		&film.PreviewPath,
		&film.MediaPath,
		&film.PreviewVideoPath,
		&film.ReleaseYear,
		&film.Rating,
		&film.AgeRestriction,
	)
	film.ID = id

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, "No rows")
		return domain.Video{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return domain.Video{}, err
	}

	return film, nil
}

func (r *filmsPostgresqlRepository) GetFilmCast(filmId int) ([]domain.Cast, error) {
	rows, err := r.db.Query(r.ctx, getFilmCastQuery, filmId)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmCast", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmCast query result:", rows)

	var artists []domain.Cast
	for rows.Next() {
		var artist domain.Cast
		err = rows.Scan(
			&artist.ID,
			&artist.Name,
			&artist.ImgPath,
		)
		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmCast", err, err.Error())
			return nil, err
		}
		artists = append(artists, artist)
	}
	if len(artists) == 0 {
		return nil, domain.ErrNotFound
	}

	return artists, nil
}

func (r *filmsPostgresqlRepository) GetCastPage(id int) ([]domain.Video, error) {
	rows, err := r.db.Query(r.ctx, getCastPageQuery, id)
	if err != nil {
		logs.LogError(logs.Logger, "cast_postgres", "GetCastPage", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetCastPage query result:", rows)

	var films []domain.Video
	for rows.Next() {
		var film domain.Video
		err = rows.Scan(
			&film.ID,
			&film.Name,
			&film.PreviewPath,
			&film.Rating,
			&film.PreviewVideoPath,
		)

		if err != nil {
			return nil, err
		}

		films = append(films, film)
	}
	if len(films) == 0 {
		return nil, domain.ErrNotFound
	}

	return films, nil
}

func (r *filmsPostgresqlRepository) GetCastName(FilmId int) (domain.Cast, error) {
	rows := r.db.QueryRow(r.ctx, getCastNameQuery, FilmId)

	logs.Logger.Debug("GetCastName query result:", rows)

	var cast domain.Cast
	err := rows.Scan(
		&cast.Name,
		&cast.Brithday,
		&cast.Place,
		&cast.Carier,
		&cast.ImgPath,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetCastName", err, err.Error())
		return domain.Cast{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetCastName", err, err.Error())
		return domain.Cast{}, err
	}

	return cast, nil
}

func (r *filmsPostgresqlRepository) GetTopRate() (domain.Video, error) {
	rows := r.db.QueryRow(r.ctx, getTopRateQuery)

	logs.Logger.Debug("GetTopRateFilm query result:", rows)

	var film domain.Video
	err := rows.Scan(
		&film.ID,
		&film.Name,
		&film.Description,
		&film.PreviewVideoPath,
		&film.MediaPath,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetTopRate", err, err.Error())
		return domain.Video{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetTopRate", err, err.Error())
		return domain.Video{}, err
	}

	return film, nil
}
