package postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"context"
)

const addAttributeQuery = `
	INSERT INTO survey (id, attribute, rate)
	VALUES ($1, $2, $3)
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
	if survey.Attribute == "" || survey.Metric == 0 {
		return domain.ErrBadRequest
	}

	result := r.db.QueryRow(r.ctx, addAttributeQuery,
		survey.
			survey.Attribute,
		survey.Metric)

	logs.Logger.Debug("AddSurvey queryRow result:", result)

	return nil
}
