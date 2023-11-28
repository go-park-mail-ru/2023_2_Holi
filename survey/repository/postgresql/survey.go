package postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"context"
)

const addAttributeQuery = `
	INSERT INTO survey (id, attribute, rate)
	VALUES ($1, $2, $3)
	ON CONFLICT (id, attribute) DO UPDATE SET rate = $3
`

const checkSurveyQuery = `
	SELECT EXISTS(SELECT 1
				  FROM survey 
				  WHERE id = $1 AND attribute = $2)
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
	result, err := r.db.Exec(r.ctx, addAttributeQuery,
		survey.ID,
		survey.Attribute,
		survey.Metric,
	)
	// if err == pgx.ErrNoRows {
	// 	logs.LogError(logs.Logger, "survey_postgres", "AddSurvey", err, err.Error())
	// 	return domain.ErrNotFound
	// }
	if err != nil {
		logs.LogError(logs.Logger, "survey_postgres", "AddSurvey", err, err.Error())
		return err
	}

	logs.Logger.Info("AddSurvey queryRow result:", result)

	return nil
}

func (r *surveyPostgresqlRepository) SurveyExists(survey domain.Survey) (bool, error) {
	result := r.db.QueryRow(r.ctx, checkSurveyQuery, survey.ID, survey.Attribute)

	var exist bool
	if err := result.Scan(&exist); err != nil {
		logs.LogError(logs.Logger, "survey_postgres", "SurveyExists", err, err.Error())
		return false, err
	}

	return exist, nil
}
