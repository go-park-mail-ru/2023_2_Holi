package movies_http

import (
	"encoding/json"
	"net/http"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type MoviesHandler struct {
	MoviesUsecase domain.MoviesUsecase
}

func NewMoviesHandler(router *mux.Router, fu domain.MoviesUsecase) {
	handler := &MoviesHandler{
		MoviesUsecase: fu,
	}

	router.HandleFunc("/v1/Moviess/genre/{genre}", handler.GetMoviesByGenre).Methods(http.MethodGet, http.MethodOptions)
}

// GetMoviessByGenre godoc
// @Summary Get Moviess by genre
// @Description Get a list of Moviess based on the specified genre.
// @Tags Moviess
// @Param genre path string true "The genre of the Moviess you want to retrieve."
// @Produce json
// @Success 200 {json} domain.Movies
// @Failure 400 {json} ApiResponse
// @Failure 404 {json} ApiResponse
// @Failure 500 {json} ApiResponse
// @Router /api/v1/Moviess/genre/{genre} [get]
func (h *MoviesHandler) GetMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]

	movies, err := h.MoviesUsecase.GetMoviesByGenre(genre)
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
			"films": movies,
		},
	}

	logs.Logger.Debug("movies:", movies)
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
