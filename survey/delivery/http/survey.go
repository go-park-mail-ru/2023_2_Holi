package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type SurveyHandler struct {
	SurveyUsecase domain.SurveyUsecase
	UtilsUsecase  domain.UtilsUsecase
}

func NewSurveyHandler(mainRouter *mux.Router, s domain.SurveyUsecase, uu domain.UtilsUsecase) {
	handler := &SurveyHandler{
		SurveyUsecase: s,
		UtilsUsecase:  uu,
	}

	mainRouter.HandleFunc("/api/v1/survey/add", handler.AddSurvey).Methods(http.MethodPost, http.MethodOptions)
}

func (s *SurveyHandler) AddSurvey(w http.ResponseWriter, r *http.Request) {

	userID := context.Get(r, "userID")

	var survey domain.Survey

	err := json.NewDecoder(r.Body).Decode(&survey)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "auth_http", "Login", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Http survey:", survey)

	//defer r.CloseAndAlert(r.Body)

	survey.Id = userID
	survey.Attribute = strings.TrimSpace(survey.Attribute)
	survey.Metric = survey.Metric

	err = s.SurveyUsecase.AddSurvey(survey)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "Add", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
