package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const getVideoByGenreQuery = `
	SELECT DISTINCT v.id, e.name, e.preview_path, v.rating , v.preview_video_path
	FROM video AS v
		JOIN video_cast AS vc ON v.id = vc.video_id
		JOIN "cast" AS c ON vc.cast_id = c.id
		JOIN episode AS e ON e.video_id = v.id
		JOIN video_genre AS vg ON v.id = vg.video_id
		JOIN genre AS g ON vg.genre_id = g.id
	WHERE g.name = $1;
`

const addToFavouritesQuery = `
	INSERT INTO  (password, name, email, image_path)
	VALUES ($1, $2, $3, $4)
	RETURNING id
`

type videoPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewVideoPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.VideoRepository {
	return &videoPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *videoPostgresqlRepository) AddToFavourites(id int) error {
	row := r.db.Exec(r.ctx, addToFavouritesQuery, id)

	logs.Logger.Debug("GetFilmData query result:", row)

	var film domain.Film
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

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, "No rows")
		return domain.Film{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return domain.Film{}, err
	}

	return film, nil
}
