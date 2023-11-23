package http

import (
	"net/http"
	"strconv"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

type FilmsHandler struct {
	FilmsUsecase domain.FilmsUsecase
}

func NewFilmsHandler(router *mux.Router, fu domain.FilmsUsecase) {
	handler := &FilmsHandler{
		FilmsUsecase: fu,
	}

	router.HandleFunc("/v1/series/genre/{genre}", handler.GetFilmsByGenre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/{id}", handler.GetFilmData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/cast/{id}", handler.GetCastPage).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/top/rate", handler.GetTopRate).Methods(http.MethodGet, http.MethodOptions)
}

// GetFilmsByGenre godoc
// @Summary 		Get films by genre
// @Description 	Get a list of films based on the specified genre.
// @Tags 			Films
// @Param 			genre path string true "The genre of the Films you want to retrieve."
// @Produce 		json
// @Success         200  {object} object{body=object{film=domain.Film}}
// @Failure         400  {object} object{err=string}
// @Failure         404  {object} object{err=string}
// @Failure         500  {object} object{err=string}
// @Router 			/api/v1/films/genre/{genre} [get]
func (h *FilmsHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]

	films, err := h.FilmsUsecase.GetFilmsByGenre(genre)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmsByGenre", err, "Failed to get films")
		return
	}

	logs.Logger.Debug("films:", films)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"films": films,
		},
		http.StatusOK,
	)
}

// GetFilmData godoc
// @Summary 		Get Film data by id
// @Description 	Get content for Film page
// @Tags 			Films
// @Param 			id path int true "Id film you want to get."
// @Produce 		json
// @Success 		200 {json} domain.Films
// @Failure 		400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
// @Router 			/api/v1/films/{id} [get]
func (h *FilmsHandler) GetFilmData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "GetFilmData", err, err.Error())
		return
	}

	film, artists, err := h.FilmsUsecase.GetFilmData(filmID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmData", err, err.Error())
		return
	}

	logs.Logger.Debug("film:", film)
	logs.Logger.Debug("artists:", artists)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"film":    film,
			"artists": artists,
		},
		http.StatusOK,
	)
}

// GetCastPage godoc
// @Summary 		Get cast page
// @Description 	Get a list of films based on the cast name.
// @Tags 			Cast
// @Param 			cast path string true "The Films of the Cast you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.Films
// @Failure			400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
// @Router 			/api/v1/films/cast/{id} [get]
func (h *FilmsHandler) GetCastPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films_http", "GetCastPage", err, err.Error())
		return
	}
	films, cast, err := h.FilmsUsecase.GetCastPage(id)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetCastPage", err, "Failed to get cast")
		return
	}

	logs.Logger.Debug("Http GetCastPage:", films)
	logs.Logger.Debug("Http GetCastPage:", cast)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"films": films,
			"cast":  cast,
		},
		http.StatusOK,
	)
}

func (h *FilmsHandler) GetTopRate(w http.ResponseWriter, r *http.Request) {
	film, err := h.FilmsUsecase.GetTopRate()
	if err != nil {
		//domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		//logs.LogError(logs.Logger, "http", "GetCastPage", err, "Failed to get cast")
		return
	}

	//logs.Logger.Debug("Http GetCastPage:", film)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"film": film,
		},
		http.StatusOK,
	)
}
