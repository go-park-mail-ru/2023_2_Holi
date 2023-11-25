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

const getStatQuery = `
	select attribute, rate, count(rate)
	from survey
	group by attribute, rate
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

func (r *surveyPostgresqlRepository) GetStat() ([]domain.Stat, error) {
	rows, err := r.db.Query(r.ctx, getStatQuery)
	if err != nil {
		logs.LogError(logs.Logger, "stat_postgres", "GetStat", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetStat query result:", rows)

	i := 0

	var stats []domain.Stat

	var cort1, cort2 string
	var cort3 int
	for rows.Next() {

		var name string

		err = rows.Scan
		cort1 = err[0]
		cort2 = err[1]
		cort3 = err[2]
		if i == 0 {
			&stats[i].Name = cort1
			a := &stats[i].Value[cort2]
			a = append(a, cort3)
			name = cort1
		}

		if name != cort1 {
			i++
		}

		a := &stats[i].Value[cort2]
		a = append(a, cort3)
	}
	logs.Logger.Info("lenth", len(stats))
	if len(stats) == 0 {
		return nil, domain.ErrNotFound
	}

	return stats, nil
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
