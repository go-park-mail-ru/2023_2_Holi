package domain

type Survey struct {
	Attribute string `json:"attribute"`
	Metric    int    `json:"metric"`
}

type SurveyUsecase interface {
	AddSurvey(survey Survey) error
}

type SurveyRepository interface {
	AddSurvey(survey Survey) error
}
