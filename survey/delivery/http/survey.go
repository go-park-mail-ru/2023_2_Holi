package http

import (
	"encoding/json"
	"errors"
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
	mainRouter.HandleFunc("/v1/survey/check/{attr}", handler.Check).Methods(http.MethodGet, http.MethodOptions)
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

	//defer r.CloseAndAlert(r.Body)

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

func (s *SurveyHandler) Check(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	attr := vars["attr"]

	if attr == "" {
		domain.WriteError(w, "empty param", http.StatusBadRequest)
		logs.LogError(logs.Logger, "survey_http", "Check", errors.New("empty param"), "Failed to decode json from body")
	}
	userID := context.Get(r, "userID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "survey_http", "Check", err, err.Error())
	}
	exist, err := s.SurveyUsecase.CheckSurvey(domain.Survey{ID: uID, Attribute: attr})
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "survey_http", "Check", err, err.Error())
	}

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"passed": exist,
		},
		http.StatusOK,
	)
}
