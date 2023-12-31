package domain

type Survey struct {
	ID        int    `json:"id"`
	Attribute string `json:"attribute"`
	Metric    string `json:"metric"`
}

type SurveyUsecase interface {
	AddSurvey(survey Survey) error
	CheckSurvey(survey Survey) (bool, error)
}

type SurveyRepository interface {
	AddSurvey(survey Survey) error
	SurveyExists(survey Survey) (bool, error)
}
