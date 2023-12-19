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
	WHERE id = $1
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
	SELECT rate
	FROM video_estimation
	WHERE video_id = $1 AND user_id = $2
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

func (r *ratingPostgresqlRepository) SelectRating(videoID int, tx *pgx.Tx) (float64, error) {
	row := (*tx).QueryRow(r.ctx, selectRatingQuery, videoID)

	logs.Logger.Debug("selectRatingQuery query result:", row)

	var rating float64
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

func (r *ratingPostgresqlRepository) Insert(rate domain.Rate) (float64, error) {
	tx, err := r.db.Begin(r.ctx)
	if err != nil {
		return 0, domain.ErrInternalServerError
	}
	defer tx.Rollback(r.ctx)

	_, err = tx.Exec(r.ctx, insertQuery, rate.Rate, rate.VideoID, rate.UserID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "Insert", err, err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, domain.ErrOutOfRange
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return 0, domain.ErrBadRequest
		}

		return 0, err
	}

	rating, err := r.SelectRating(rate.VideoID, &tx)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return 0, err
	}

	return rating, nil
}

func (r *ratingPostgresqlRepository) Delete(rate domain.Rate) (float64, error) {
	tx, err := r.db.Begin(r.ctx)
	if err != nil {
		return 0, domain.ErrInternalServerError
	}
	defer tx.Rollback(r.ctx)

	res, err := tx.Exec(r.ctx, deleteQuery, rate.VideoID, rate.UserID)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "Delete", err, err.Error())
		return 0, err
	}

	if res.RowsAffected() == 0 {
		logs.LogError(logs.Logger, "postgresql(rating)", "Delete", domain.ErrOutOfRange, domain.ErrOutOfRange.Error())
		return 0, domain.ErrOutOfRange
	}

	rating, err := r.SelectRating(rate.VideoID, &tx)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return 0, err
	}

	return rating, nil
}

func (r *ratingPostgresqlRepository) Exists(rate domain.Rate) (bool, int, error) {
	result := r.db.QueryRow(r.ctx, existsQuery, rate.VideoID, rate.UserID)

	var rateNumber int
	err := result.Scan(&rateNumber)

	if err == pgx.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		logs.LogError(logs.Logger, "postgresql(rating)", "Exists", err, err.Error())
		return false, 0, err
	}

	return true, rateNumber, nil
}
