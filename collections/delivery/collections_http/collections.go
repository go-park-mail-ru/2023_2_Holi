package collections_http

import (
	"encoding/json"
	"net/http"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logs"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type FilmHandler struct {
	FilmUsecase domain.FilmUsecase
}

func NewFilmHandler(router *mux.Router, fu domain.FilmUsecase) {
	handler := &FilmHandler{
		FilmUsecase: fu,
	}

	router.HandleFunc("/v1/films/genre/{genre}", handler.GetFilmsByGenre).Methods(http.MethodGet, http.MethodOptions)
}

// GetFilmsByGenre godoc
// @Summary Get films by genre
// @Description Get a list of films based on the specified genre.
// @Tags films
// @Param genre path string true "The genre of the films you want to retrieve."
// @Produce json
// @Success 200 {array} Film
// @Failure 400 {application/json} ErrorResponse
// @Failure 404 {application/json} ErrorResponse
// @Failure 500 {application/json} ErrorResponse
// @Router /api/v1/films/genre/{genre} [get]

func (h *FilmHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]

	films, err := h.FilmUsecase.GetFilmsByGenre(genre)
	if err != nil {
		response := ApiResponse{
			Status: getStatusCode(err),
			Body: map[string]string{
				"error": err.Error(),
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := ApiResponse{
		Status: http.StatusOK,
		Body: map[string]interface{}{
			"films": films,
		},
	}

	logs.Logger.Debug("Films:", films)
	json.NewEncoder(w).Encode(response)
}

func (h *FilmHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	imagePath := "static/preview_path/" + filename

	http.ServeFile(w, r, imagePath)
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
