package http

import (
	"encoding/json"
	"net/http"

	"2023_2_Holi/domain"

	logs "2023_2_Holi/logs"

	"github.com/gorilla/mux"
)

var logger = logs.LoggerInit()

type ApiResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type ArtistHandler struct {
	ArtistUsecase domain.ArtistUsecase
}

func NewArtistHandler(router *mux.Router, au domain.ArtistUsecase) {
	handler := &ArtistHandler{
		ArtistUsecase: au,
	}

	router.HandleFunc("/v1/artist/{name}/{surname}", handler.GetArtistPage).Methods(http.MethodGet, http.MethodOptions)
}

func (h *ArtistHandler) GetArtistPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	surname := vars["surname"]

	films, err := h.ArtistUsecase.GetArtistPage(name, surname)
	if err != nil {
		response := ApiResponse{
			Status: getStatusCode(err),
			Body: map[string]string{
				"error": err.Error(),
			},
		}
		logs.LogError(logs.Logger, "http", "GetArtistPage", err, "Failed to get artist")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := ApiResponse{
		Status: http.StatusOK,
		Body: map[string]interface{}{
			"films": films,
		},
	}

	logger.Debug("Http GetArtistPage:", films)
	json.NewEncoder(w).Encode(response)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrWrongCredentials:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
