package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"
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

	mainRouter.HandleFunc("/v1/survey/add", handler.AddSurvey).Methods(http.MethodPost, http.MethodOptions)
	mainRouter.HandleFunc("/v1/survey/stat", handler.AddSurvey).Methods(http.MethodGet, http.MethodOptions)

}

func (s *SurveyHandler) AddSurvey(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		return
	}

	var survey domain.Survey

	err = json.NewDecoder(r.Body).Decode(&survey)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "auth_http", "Login", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Http survey:", survey)

	survey.ID = userID
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

func (s *SurveyHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stat, err := s.SurveyUsecase.GetStat()
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetGenres", err, err.Error())
		return
	}

	logs.Logger.Debug("Http GetGenres:", stat)
	domain.WriteResponse(
		w,
		http.StatusOK,
	)
}
