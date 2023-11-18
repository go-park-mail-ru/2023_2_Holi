package http

import (
	"net/http"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

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
// @Success 		200 {array} domain.Genre
// @Failure			400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
// @Router /v1/genres [get]
func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.GenreUsecase.GetGenres()
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetGenres", err, err.Error())
		return
	}

	logs.Logger.Debug("Http GetGenres:", genres)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"genres": genres,
		},
		http.StatusOK,
	)
}
