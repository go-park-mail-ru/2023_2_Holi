package postgresql

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

const subQuery = `
	UPDATE "user_subscription"
		SET subscription_up_to = $1
		WHERE user_id = $2;
`
const subDeleteQuery = `
	UPDATE "user_subscription"
		SET subscription_up_to = '01-01-0001'
		WHERE user_id = $1;
`

const checkSubQuery = `
	SELECT subscription_up_to
		FROM "user_subscription"
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

func (r *subsPostgresqlRepository) Subscribe(subId int) error {
	subTime := time.Now()

	subTime = subTime.AddDate(0, 1, 0)

	logs.Logger.Debug("until sub time:", subTime)

	row, err := r.db.Exec(r.ctx, subQuery, subTime, subId)

	if err != nil {
		logs.LogError(logs.Logger, "subc_postgres", "GetSubData", err, err.Error())
		return err
	}

	logs.Logger.Debug("subscribe query result:", row)

	return nil
}

func (r *subsPostgresqlRepository) UnSubscribe(subId int) error {
	row, err := r.db.Exec(r.ctx, subDeleteQuery, subId)
	if err != nil {
		logs.LogError(logs.Logger, "subs_postgresql", "UnSubscribe", err, err.Error())
		return err
	}
	logs.Logger.Debug("UnSubscribe query result:", row)

	return nil
}

func (r *subsPostgresqlRepository) CheckSub(subId int) (sub time.Time, error error) {
	row := r.db.QueryRow(r.ctx, checkSubQuery, subId)

	logs.Logger.Debug("GetSubData query result:", row)

	var subUpTo *time.Time

	err := row.Scan(&subUpTo)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "subc_postgres", "GetSubData", err, err.Error())
		return *subUpTo, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "subc_postgres", "GetSubData", err, err.Error())
		return *subUpTo, err
	}

	return *subUpTo, err
}
