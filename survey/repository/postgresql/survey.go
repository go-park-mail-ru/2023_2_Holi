package postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"context"
)

const addAttributeQuery = `
	INSERT INTO survey (id, attribute, rate)
	VALUES ($1, $2, $3)
	ON CONFLICT DO UPDATE SET rate = $3
`

const getStatQuery = `

`

type surveyPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewSurveyPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.SurveyRepository {
	return &surveyPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *surveyPostgresqlRepository) AddSurvey(survey domain.Survey) error {

	result := r.db.QueryRow(r.ctx, addAttributeQuery,
		survey.Attribute,
		survey.Metric,
		survey.ID)

	logs.Logger.Debug("AddSurvey queryRow result:", result)

	return nil
}

func (r *surveyPostgresqlRepository) GetStat() ([]domain.Stat, error) {
	rows, err := r.db.Query(r.ctx, getStatQuery)
	if err != nil {
		logs.LogError(logs.Logger, "stat_postgres", "GetStat", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetStat query result:", rows)

	var stats []domain.Stat
	for rows.Next() {
		var stat domain.Genre
		err = rows.Scan()

		if err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}
	logs.Logger.Info("lenth", len(stats))
	if len(stats) == 0 {
		return nil, domain.ErrNotFound
	}

	return stats, nil
}
