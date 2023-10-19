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

type GenreHandler struct {
	GenreUsecase domain.GenreUsecase
}

func NewGenreHandler(router *mux.Router, gu domain.GenreUsecase) {
	handler := &GenreHandler{
		GenreUsecase: gu,
	}

	router.HandleFunc("/v1/genres", handler.GetGenres).Methods(http.MethodGet, http.MethodOptions)
}

// GetGenres godoc
// @Summary Get genres
// @Description Get a list of genres.
// @Tags genres
// @Produce json
// @Success 200 {array} Genre
// @Failure 500 {object} ErrorResponse
// @Router /v1/genres [get]

func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {

	genres, err := h.GenreUsecase.GetGenres()
	if err != nil {
		response := ApiResponse{
			Status: getStatusCode(err),
			Body: map[string]string{
				"error": err.Error(),
			},
		}
		logs.LogError(logs.Logger, "http", "GetGenres", err, "Failed to get genres")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := ApiResponse{
		Status: http.StatusOK,
		Body: map[string]interface{}{
			"genres": genres,
		},
	}

	logger.Debug("Http GetGenres:", genres)
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
