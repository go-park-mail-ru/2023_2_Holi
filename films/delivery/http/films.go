package films_http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type FilmsHandler struct {
	FilmsUsecase domain.FilmsUsecase
}

func NewFilmsHandler(router *mux.Router, fu domain.FilmsUsecase) {
	handler := &FilmsHandler{
		FilmsUsecase: fu,
	}

	router.HandleFunc("/v1/films/genre/{genre}", handler.GetFilmsByGenre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/films/{id}", handler.GetFilmData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/films/cast/{id}", handler.GetCastPage).Methods(http.MethodGet, http.MethodOptions)
}

// GetFilmsByGenre godoc
// @Summary 		Get films by genre
// @Description 	Get a list of films based on the specified genre.
// @Tags 			Films
// @Param 			genre path string true "The genre of the Films you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.Films
// @Failure			400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router 			/api/v1/films/genre/{genre} [get]
func (h *FilmsHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]

	Films, err := h.FilmsUsecase.GetFilmsByGenre(genre)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmsByGenre", err, "Failed to get films")
		return
	}
	response := map[string]interface{}{
		"films": Films,
	}

	logs.Logger.Debug("Films:", Films)
	json.NewEncoder(w).Encode(response)
}

// GetFilmData godoc
// @Summary 		Get Film data by id
// @Description 	Get content for Film page
// @Tags 			Films
// @Param 			id path int true "The genre of the Films you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.Films
// @Failure 		400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router 			/api/v1/Films/{id} [get]
func (h *FilmsHandler) GetFilmData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	FilmID := vars["id"]
	id, err := strconv.Atoi(FilmID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "Films_http", "GetFilmData", err, err.Error())
		return
	}

	Film, artists, err := h.FilmsUsecase.GetFilmData(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "Films_http", "GetFilmData", err, err.Error())
		return
	}

	response := map[string]interface{}{
		"Film":    Film,
		"artists": artists,
	}

	logs.Logger.Debug("Film:", Film)
	logs.Logger.Debug("artists:", artists)
	json.NewEncoder(w).Encode(response)
}

// GetCastPage godoc
// @Summary 		Get cast page
// @Description 	Get a list of films based on the cast name.
// @Tags 			Cast
// @Param 			cast path string true "The Films of the Cast you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.Films
// @Failure			400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router 			api/v1/films/cast/{id} [get]
func (h *FilmsHandler) GetCastPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	CastID := vars["id"]
	id, err := strconv.Atoi(CastID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "Films_http", "GetCastPage", err, err.Error())
		return
	}
	films, cast, err := h.FilmsUsecase.GetCastPage(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetCastPage", err, "Failed to get cast")
		return
	}
	response := map[string]interface{}{
		"films": films,
		"cast":  cast,
	}

	logs.Logger.Debug("Http GetArtistPage:", films)
	logs.Logger.Debug("Http GetArtistPage:", cast)
	json.NewEncoder(w).Encode(response)
}
