package postgresql

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"time"

	"context"
)

const subQuery = `
    INSERT INTO user_subscription (user_id, subscription_start, subscription_up_to)
	VALUES ($1, $2, $3)
`
const subDeleteQuery = `
	UPDATE user_subscription
	SET subscription_start = NULL, subscription_up_to = NULL
	WHERE user_id = $1;
`

const checkSubQuery = `
	SELECT subscription_start, subscription_up_to
	FROM user_subscription
	WHERE user_id = $1;
`

const checkExistQuery = `
	INSERT INTO user_subscription (user_id, subscription_start, subscription_up_to)
	VALUES ($1, $2, $3)
`

const getInfoSubQuery = `
	SELECT subscription_start, subscription_up_to, user_id
	FROM user_subscription
	WHERE user_id = $1;
`

type subsPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewSubsPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.SubsRepository {
	return &subsPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *subsPostgresqlRepository) Subscribe(subId int, flag int) error {
	currentTime := time.Now()
	subTime := time.Now()
	if flag == 1 {
		subTime = currentTime.AddDate(0, 1, 0)
	}
	if flag == 6 {
		subTime = currentTime.AddDate(0, 6, 0)
	}
	rows := r.db.QueryRow(r.ctx, subQuery, subId, currentTime, subTime)

	logs.Logger.Debug("Subscribe query result:", rows)

	var id int
	if err := rows.Scan(&id); err != nil {
		logs.LogError(logs.Logger, "subs_postgres", "Subscribe", err, err.Error())
		return err
	}
	return nil
}

func (r *subsPostgresqlRepository) UnSubscribe(subId int) error {
	rows, err := r.db.Query(r.ctx, subQuery, subId)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmsByGenre query result:", rows)

	return nil
}

func (r *subsPostgresqlRepository) CheckSub(subId int) error {
	rows, err := r.db.Query(r.ctx, subQuery, subId)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmsByGenre query result:", rows)

	return nil
}

func (r *subsPostgresqlRepository) GetSubInfo(subId int) (domain.SubInfo, error) {
	rows, err := r.db.Query(r.ctx, subQuery, subId)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return domain.SubInfo{}, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetFilmsByGenre query result:", rows)

	return domain.SubInfo{}, nil
}
