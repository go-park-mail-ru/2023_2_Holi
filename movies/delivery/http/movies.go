package movies_http

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	router.HandleFunc("/v1/movies/genre/{genre}", handler.GetMoviesByGenre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/movies/{id}", handler.GetMovieData).Methods(http.MethodGet, http.MethodOptions)
}

// GetMoviesByGenre godoc
// @Summary 		Get movies by genre
// @Description 	Get a list of movies based on the specified genre.
// @Tags 			movies
// @Param 			genre path string true "The genre of the movies you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.movies
// @Failure			400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router 			/api/v1/movies/genre/{genre} [get]
func (h *MoviesHandler) GetMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]

	movies, err := h.MoviesUsecase.GetMoviesByGenre(genre)
	if err != nil {
		response := ApiResponse{
			Status: domain.GetStatusCode(err),
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

// GetMovieData godoc
// @Summary 		Get movie data by id
// @Description 	Get content for movie page
// @Tags 			movies
// @Param 			id path int true "The genre of the movies you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.Movies
// @Failure 		400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router 			/api/v1/movies/{id} [get]
func (h *MoviesHandler) GetMovieData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["id"]
	id, err := strconv.Atoi(movieID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "movies_http", "GetMovieData", err, err.Error())
		return
	}

	movie, artists, err := h.MoviesUsecase.GetMovieData(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "movies_http", "GetMovieData", err, err.Error())
		return
	}

	response := ApiResponse{
		Status: http.StatusOK,
		Body: map[string]interface{}{
			"movie":   movie,
			"artists": artists,
		},
	}

	logs.Logger.Debug("movie:", movie)
	logs.Logger.Debug("artists:", artists)
	json.NewEncoder(w).Encode(response)
}
