package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type SurveyHandler struct {
	SurveyUsecase domain.SurveyUsecase
}

func NewSurveyHandler(mainRouter *mux.Router, s domain.SurveyUsecase) {
	handler := &SurveyHandler{
		SurveyUsecase: s,
	}

	mainRouter.HandleFunc("/api/v1/survey/add", handler.AddSurvey).Methods(http.MethodPost, http.MethodOptions)
}

func (s *SurveyHandler) AddSurvey(w http.ResponseWriter, r *http.Request) {

	var survey domain.Survey

	err := json.NewDecoder(r.Body).Decode(&survey)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "auth_http", "Login", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Http survey:", survey)

	//defer r.CloseAndAlert(r.Body)

	survey.Attribute = strings.TrimSpace(survey.Attribute)
	survey.Metric = survey.Metric

	domain.WriteResponse(
		w,
		nil,
		http.StatusOK,
	)
}
