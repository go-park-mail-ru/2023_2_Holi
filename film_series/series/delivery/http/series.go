package http

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

type SeriesHandler struct {
	SeriesUsecase domain.SeriesUsecase
}

func NewSeriesHandler(router *mux.Router, su domain.SeriesUsecase) {
	handler := &SeriesHandler{
		SeriesUsecase: su,
	}

	router.HandleFunc("/v1/series/genre/{genre}", handler.GetSeriesByGenre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/{id}", handler.GetSeriesData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/cast/{id}", handler.GetCastPageSeries).Methods(http.MethodGet, http.MethodOptions)
	// router.HandleFunc("/v1/series/top/rate", handler.GetTopRate).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/metrics", promhttp.Handler())
}

// GetSeriesByGenre godoc
// @Summary 		Get series by genre
// @Description 	Get a list of series based on the specified genre.
// @Tags 			Series
// @Param 			genre path string true "The genre of the Series you want to retrieve."
// @Produce 		json
// @Success         200  {object} object{body=object{film=domain.Video}}
// @Failure         400  {object} domain.Response
// @Failure         404  {object} domain.Response
// @Failure         500  {object} domain.Response
// @Router 			/api/v1/series/genre/{genre} [get]
func (h *SeriesHandler) GetSeriesByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre := vars["genre"]

	films, err := h.SeriesUsecase.GetSeriesByGenre(genre)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
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

// GetSeriesData godoc
// @Summary 		Get Series data by id
// @Description 	Get content for Series page
// @Tags 			Series
// @Param 			id path int true "Id series you want to get."
// @Produce 		json
// @Success 		200 {json} domain.Video, []domain.Cast, []domain.Episode
// @Failure 		400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
// @Router 			/api/v1/series/{id} [get]
func (h *SeriesHandler) GetSeriesData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "GetFilmData", err, err.Error())
		return
	}

	series, artists, episodes, err := h.SeriesUsecase.GetSeriesData(filmID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmData", err, err.Error())
		return
	}

	logs.Logger.Debug("film:", series)
	logs.Logger.Debug("artists:", artists)
	logs.Logger.Debug("episodes:", episodes)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"film":     series,
			"artists":  artists,
			"episodes": episodes,
		},
		http.StatusOK,
	)
}

// GetCastPageSeries godoc
// @Summary 		Get cast page series
// @Description 	Get a list of series based on the cast name.
// @Tags 			Series
// @Param 			cast path string true "The Series of the Cast you want to retrieve."
// @Produce 		json
// @Success 		200 {json} []domain.Video
// @Failure			400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
// @Router 			/api/v1/series/cast/{id} [get]
func (h *SeriesHandler) GetCastPageSeries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films_http", "GetCastPageSeries", err, err.Error())
		return
	}
	series, cast, err := h.SeriesUsecase.GetCastPageSeries(id)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetCastPageSeries", err, "Failed to get cast")
		return
	}

	logs.Logger.Debug("Http GetCastPageSeries:", series)
	logs.Logger.Debug("Http GetCastPageSeries:", cast)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"series": series,
			"cast":   cast,
		},
		http.StatusOK,
	)
}

// func (h *FilmsHandler) GetTopRate(w http.ResponseWriter, r *http.Request) {
// 	film, err := h.FilmsUsecase.GetTopRate()
// 	if err != nil {
// 		//domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
// 		//logs.LogError(logs.Logger, "http", "GetCastPage", err, "Failed to get cast")
// 		return
// 	}

// 	//logs.Logger.Debug("Http GetCastPage:", film)
// 	domain.WriteResponse(
// 		w,
// 		map[string]interface{}{
// 			"film": film,
// 		},
// 		http.StatusOK,
// 	)
// }
