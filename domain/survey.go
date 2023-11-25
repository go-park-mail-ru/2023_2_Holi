package domain

type Survey struct {
	ID        int    `json:"id"`
	Attribute string `json:"attribute"`
	Metric    string `json:"metric"`
}

type Stat struct {
	Name  string `json:"name"`
	Value map[string]interface{}
}

type SurveyUsecase interface {
	AddSurvey(survey Survey) error
	GetStat() ([]Stat, error)
}

type SurveyRepository interface {
	AddSurvey(survey Survey) error
	GetStat() ([]Stat, error)
}
