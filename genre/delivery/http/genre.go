package http

import (
	"encoding/json"
	"net/http"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var logger = logs.LoggerInit()

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
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
// @Success 		200 {array} Genre
// @Failure			400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router /v1/genres [get]

func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {

	genres, err := h.GenreUsecase.GetGenres()
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetGenres", err, "Failed to get genres")
		return
	}
	response := map[string]interface{}{
		"genres": genres,
	}

	logger.Debug("Http GetGenres:", genres)
	json.NewEncoder(w).Encode(&Result{Body: response})
}
