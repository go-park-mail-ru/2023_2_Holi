package http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type SubsHandler struct {
	SubsUsecase domain.SubsUsecase
}

func NewSubsHandler(router *mux.Router, su domain.SubsUsecase) {
	handler := &SubsHandler{
		SubsUsecase: su,
	}

	router.HandleFunc("/v1/subs/sub/{id}/{flag}", handler.Subscribe).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/subs/unsub/{id}", handler.UnSubscribe).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/v1/subs/check/{id}", handler.CheckSub).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/subs/info/{id}", handler.GetSubInfo).Methods(http.MethodGet, http.MethodOptions)
}

func (h *SubsHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "Subscribe", err, err.Error())
		return
	}
	flag, err := strconv.Atoi(vars["flag"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "Subscribe", err, err.Error())
		return
	}

	err = h.SubsUsecase.Subscribe(subID, flag)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "subs_http", "Subscribe", err, "Failed to sub")
		return
	}

	logs.Logger.Debug("subs:", err)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"responce": "successful Subscription",
		},
		http.StatusOK,
	)
}

func (h *SubsHandler) UnSubscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "UnSubscribe", err, err.Error())
		return
	}

	err = h.SubsUsecase.UnSubscribe(subID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "subs_http", "UnSubscribe", err, "Failed to Unsub")
		return
	}

	logs.Logger.Debug("subs:", err)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"responce": "successful Unsubscription",
		},
		http.StatusOK,
	)
}

func (h *SubsHandler) CheckSub(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "CheckSub", err, err.Error())
		return
	}

	err = h.SubsUsecase.UnSubscribe(subID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "subs_http", "CheckSub", err, "Failed to CheckSub")
		return
	}

	logs.Logger.Debug("subs:", err)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"responce": "user is sub",
		},
		http.StatusOK,
	)
}

func (h *SubsHandler) GetSubInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "CheckSub", err, err.Error())
		return
	}

	suber, err := h.SubsUsecase.GetSubInfo(subID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "subs_http", "CheckSub", err, "Failed to CheckSub")
		return
	}

	logs.Logger.Debug("suber:", suber)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"suber": suber,
		},
		http.StatusOK,
	)
}
