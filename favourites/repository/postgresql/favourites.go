package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const addToFavouritesQuery = `
	INSERT INTO favourite (video_id, user_id)
	VALUES ($1, $2)
`

const DeleteFromFavouritesQuery = `
	DELETE FROM favourite
	WHERE video_id = $1 AND user_id = $2
`

const GetAllSeriesQuery = `
	SELECT v.id, v.name, v.description,
		v.preview_path, v.preview_video_path, v.release_year, v.rating, v.age_restriction, v.seasons_count
	FROM video v
		JOIN favourite f ON video_id = v.id
	WHERE f.user_id = $1 AND v.seasons_count != 0
`

const GetAllFilmsQuery = `
	SELECT v.id, e.name, e.description, e.duration,
		e.preview_path, e.media_path, v.preview_video_path, v.release_year, v.rating, v.age_restriction
	FROM video v
		JOIN episode e ON e.video_id = v.id
		JOIN favourite f ON f.video_id = v.id
	WHERE f.user_id = $1 AND v.seasons_count = 0
`

type favouritesUsecasePostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewFavouritesPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.FavouritesRepository {
	return &favouritesUsecasePostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *favouritesUsecasePostgresqlRepository) Add(videoID, userID int) error {
	_, err := r.db.Exec(r.ctx, addToFavouritesQuery, videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "Add", err, err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.ErrOutOfRange
		}

		return err
	}

	return nil
}

func (r *favouritesUsecasePostgresqlRepository) Delete(videoID, userID int) error {
	res, err := r.db.Exec(r.ctx, DeleteFromFavouritesQuery, videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "Delete", err, err.Error())
		return err
	}

	if res.RowsAffected() == 0 {
		logs.LogError(logs.Logger, "postgresql", "Delete", domain.ErrOutOfRange, domain.ErrOutOfRange.Error())
		return domain.ErrOutOfRange
	}

	return nil
}

func (r *favouritesUsecasePostgresqlRepository) GetAll(userID int) ([]domain.Video, error) {
	rows, err := r.db.Query(r.ctx, GetAllSeriesQuery, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "GetAll", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetAllSeries query result:", rows)

	var videos []domain.Video

	for rows.Next() {
		var series domain.Video
		err = rows.Scan(
			&series.ID,
			&series.Name,
			&series.Description,
			&series.PreviewPath,
			&series.PreviewVideoPath,
			&series.ReleaseYear,
			&series.Rating,
			&series.AgeRestriction,
			&series.SeasonsCount,
		)

		if err != nil {
			return nil, err
		}

		videos = append(videos, series)
	}
	logs.Logger.Debug("GetAllFilmsQuery videos:", videos)

	rows, err = r.db.Query(r.ctx, GetAllFilmsQuery, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "GetAll", err, err.Error())
		return nil, err
	}
	logs.Logger.Debug("GetAllFilms query result:", rows)

	for rows.Next() {
		var film domain.Video
		err = rows.Scan(
			&film.ID,
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

		if err != nil {
			return nil, err
		}

		videos = append(videos, film)
	}
	logs.Logger.Debug("GetAllFilmsQuery videos:", videos)

	return videos, nil
}
