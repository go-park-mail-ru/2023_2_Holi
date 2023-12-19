package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const selectRatingQuery = `
	SELECT rating
	FROM video
	WHERE video_id = $1
`

const insertQuery = `
	INSERT INTO video_estimation(rate, video_id, user_id)
	VALUES ($1, $2, $3)
	ON CONFLICT(video_id, user_id)
	DO UPDATE
	SET rate = EXCLUDED.rate
`

const deleteQuery = `
	DELETE FROM video_estimation
	WHERE video_id = $1 AND user_id = $2
`

const existsQuery = `
	SELECT EXISTS(SELECT 1
				  FROM video_estimation
				  WHERE video_id = $1 AND user_id = $2)
`

type ratingPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewRatingPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.RatingRepository {
	return &ratingPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *ratingPostgresqlRepository) SelectRating(videoID int) (int, error) {
	row := r.db.QueryRow(r.ctx, selectRatingQuery, videoID)

	logs.Logger.Debug("selectRatingQuery query result:", row)

	var rating int
	err := row.Scan(
		&rating,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "postgresql(rating)", "SelectRating", err, "No rows")
		return 0, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "SelectRating", err, err.Error())
		return 0, err
	}

	return rating, nil
}

func (r *ratingPostgresqlRepository) Insert(rate domain.Rate) error {
	_, err := r.db.Exec(r.ctx, insertQuery, rate.Rate, rate.VideoID, rate.UserID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "Insert", err, err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.ErrOutOfRange
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return domain.ErrBadRequest
		}

		return err
	}

	return nil
}

func (r *ratingPostgresqlRepository) Delete(rate domain.Rate) error {
	res, err := r.db.Exec(r.ctx, deleteQuery, rate.VideoID, rate.UserID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "Delete", err, err.Error())
		return err
	}

	if res.RowsAffected() == 0 {
		logs.LogError(logs.Logger, "postgresql(rating)", "Delete", domain.ErrOutOfRange, domain.ErrOutOfRange.Error())
		return domain.ErrOutOfRange
	}

	return nil
}

func (r *ratingPostgresqlRepository) Exists(rate domain.Rate) (bool, error) {
	result := r.db.QueryRow(r.ctx, existsQuery, rate.VideoID, rate.UserID)

	var exist bool
	if err := result.Scan(&exist); err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "Exists", err, err.Error())
		return false, err
	}

	return exist, nil
}
