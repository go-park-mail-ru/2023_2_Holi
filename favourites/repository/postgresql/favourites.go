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

const SelectAllQuery = `
	SELECT v.id, v.name, v.description,
		v.preview_path, v.preview_video_path, v.release_year, v.rating, v.age_restriction, v.seasons_count
	FROM video AS v
		JOIN favourite AS f ON video_id = v.id
	WHERE f.user_id = $1
`
const favouriteExistsQuery = `
	SELECT EXISTS(SELECT 1
				  FROM favourite
				  WHERE video_id = $1 AND user_id = $2)
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

func (r *favouritesUsecasePostgresqlRepository) InsertIntoFavourites(videoID, userID int) error {
	_, err := r.db.Exec(r.ctx, addToFavouritesQuery, videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "InsertIntoFavourites", err, err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.ErrOutOfRange
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r *favouritesUsecasePostgresqlRepository) DeleteFromFavourites(videoID, userID int) error {
	res, err := r.db.Exec(r.ctx, DeleteFromFavouritesQuery, videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "DeleteFromFavourites", err, err.Error())
		return err
	}

	if res.RowsAffected() == 0 {
		logs.LogError(logs.Logger, "postgresql", "DeleteFromFavourites", domain.ErrOutOfRange, domain.ErrOutOfRange.Error())
		return domain.ErrOutOfRange
	}

	return nil
}

func (r *favouritesUsecasePostgresqlRepository) SelectAllFavourites(userID int) ([]domain.Video, error) {
	rows, err := r.db.Query(r.ctx, SelectAllQuery, userID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "SelectAllFavourites", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("SelectAllFavourites query result:", rows)

	var videos []domain.Video

	for rows.Next() {
		var video domain.Video
		err = rows.Scan(
			&video.ID,
			&video.Name,
			&video.Description,
			&video.PreviewPath,
			&video.PreviewVideoPath,
			&video.ReleaseYear,
			&video.Rating,
			&video.AgeRestriction,
			&video.SeasonsCount,
		)

		if err != nil {
			return nil, err
		}

		videos = append(videos, video)
	}
	logs.Logger.Debug("SelectAllFavourites videos:", videos)

	return videos, nil
}

func (r *favouritesUsecasePostgresqlRepository) Exists(videoID, userID int) (bool, error) {
	result := r.db.QueryRow(r.ctx, favouriteExistsQuery, videoID, userID)

	var exist bool
	if err := result.Scan(&exist); err != nil {
		logs.LogError(logs.Logger, "postgresql", "Exists", err, err.Error())
		return false, err
	}

	return exist, nil
}
