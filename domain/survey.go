package domain

type Survey struct {
	ID        int    `json:"id"`
	Attribute string `json:"attribute"`
	Metric    string `json:"metric"`
}

type Stat struct {
	Name  string `json:"name"`
	Value map[string][]string
}

type SurveyUsecase interface {
	AddSurvey(survey Survey) error
	GetStat() ([]Stat, error)
	CheckSurvey(survey Survey) (bool, error)
}

type SurveyRepository interface {
	AddSurvey(survey Survey) error
	GetStat() ([]Stat, error)
	SurveyExists(survey Survey) (bool, error)
}
