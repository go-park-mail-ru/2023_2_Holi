package http

import (
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

	router.HandleFunc("/v1/series/genre/{id}", handler.GetSeriesByGenre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/{id}", handler.GetSeriesData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/cast/{id}", handler.GetCastPageSeries).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/series/top/rate", handler.GetTopRate).Methods(http.MethodGet, http.MethodOptions)
}

// GetSeriesByGenre godoc
// @Summary 		Get series by genre
// @Description 	Get a list of series based on the specified genre.
// @Tags 			Series
// @Param 			genre path string true "The Series of the genre you want to retrieve."
// @Produce 		json
//
//	@Success		200		{json}	object{body=object{[]domain.Video}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//
// @Router 			/api/v1/series/genre/{genreId} [get]
func (h *SeriesHandler) GetSeriesByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films_http", "GetCastPage", err, err.Error())
		return
	}

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
//
//	@Success		200		{json}	object{body=object{{domain.Video, []domain.Cast, []domain.Episode}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//
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
//
//	@Success		200		{json}	object{body=object{[]domain.Video}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//
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

// GetTopRate godoc
//
//	@Summary		Get top rate information
//	@Description	Get information about the most rated series.
//	@Tags			Series
//	@Param			rate	path	string	true	"The top rate Series  you want to retrieve."
//	@Produce		json
//
//	@Success		200		{json}	object{body=object{[]domain.Video}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//
//	@Router			/api/v1/series/top/rate [get]
func (h *SeriesHandler) GetTopRate(w http.ResponseWriter, r *http.Request) {
	topRate, err := h.SeriesUsecase.GetTopRate()
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "series_http", "GetTopRate", err, "Failed to get top rate series")
		return
	}

	logs.Logger.Debug("series_http GetTopRate:", topRate)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"series": topRate,
		},
		http.StatusOK,
	)
}
