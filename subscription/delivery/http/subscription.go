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

	router.HandleFunc("/v1/subs/take_request", handler.TakeRequest).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/subs/pay/{id}", handler.Pay).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/subs/sub/{id}", handler.Subscribe).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/subs/unsub/{id}", handler.UnSubscribe).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/subs/check/{id}", handler.CheckSub).Methods(http.MethodGet, http.MethodOptions)
}

func (h *SubsHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "Subscribe", err, err.Error())
		return
	}

	err = h.SubsUsecase.Subscribe(subID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "subs_http", "Subscribe", err, "Failed to sub")
		return
	}

	logs.Logger.Debug("subs:", err)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"responce": "sucsesfull sibscribe",
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

	subUpTo, status, err := h.SubsUsecase.CheckSub(subID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "subs_http", "CheckSub", err, "Failed to CheckSub")
		return
	}

	logs.Logger.Debug("subs:", err)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"subUpTo": subUpTo,
			"status":  status,
		},
		http.StatusOK,
	)
}

func (h *SubsHandler) Pay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "subs_http", "CheckSub", err, err.Error())
		return
	}

	payment := domain.Payment(userId)
	logs.Logger.Debug("payment:", payment)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"payment": payment,
		},
		http.StatusOK,
	)

}

func (h *SubsHandler) TakeRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	receivedHash := r.FormValue("sha1_hash")
	if receivedHash == "" {
		http.Error(w, "sha1_hash not provided", http.StatusBadRequest)
		return
	}

	parametersString := domain.CreateParametersString(r)

	sha1Hash, err := domain.CalculateSHA1Hash(parametersString)
	if err != nil {
		http.Error(w, "Error calculating SHA-1 hash", http.StatusInternalServerError)
		return
	}

	if sha1Hash != receivedHash {
		http.Error(w, "Invalid sha1_hash", http.StatusForbidden)
		return
	}
	_, err = http.Post("http://localhost:3006/api/v1/subs/sub/1", "application/json", nil)
	if err != nil {
		http.Error(w, "Error sending HTTP request", http.StatusInternalServerError)
		return
	}
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"status": "successful",
		},
		http.StatusOK,
	)
}
