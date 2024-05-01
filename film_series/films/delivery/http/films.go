package http

import (
	"2023_2_Holi/domain/grpc/subscription"
	"context"
	gorilla_ctx "github.com/gorilla/context"
	"net/http"
	"strconv"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

type FilmsHandler struct {
	FilmsUsecase domain.FilmsUsecase
	SubClient    subscription.SubCheckerClient
}

func NewFilmsHandler(router *mux.Router, fu domain.FilmsUsecase, subClient subscription.SubCheckerClient) {
	handler := &FilmsHandler{
		FilmsUsecase: fu,
		SubClient:    subClient,
	}

	router.HandleFunc("/v1/films/genre/{id}", handler.GetFilmsByGenre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/films/{id}", handler.GetFilmData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/films/cast/{id}", handler.GetCastPage).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/films/top/rate", handler.GetTopRate).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/films/recommendations/{id}", handler.GetRecommendations).Methods(http.MethodGet, http.MethodOptions)
}

// GetFilmsByGenre godoc
//
//	@Summary		Get films by genre
//	@Description	Get a list of films based on the specified genre.
//	@Tags			Films
//	@Param			genre	path	string	true	"The Films of the genre you want to retrieve."
//	@Produce		json
//
//	@Success		200		{json}	object{body=object{[]domain.Video}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//
//	@Router			/api/v1/films/genre/{genreId} [get]
func (h *FilmsHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genre, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films_http", "GetCastPage", err, err.Error())
		return
	}

	films, err := h.FilmsUsecase.GetFilmsByGenre(genre)
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

// GetFilmData godoc
//
//	@Summary		Get Video data by id
//	@Description	Get content for Video page
//	@Tags			Films
//	@Param			id	path	int	true	"Id film you want to get."
//	@Produce		json
//
// @Success 		200 {json} domain.Video, []domain.Cast
// @Failure 		400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
//
// @Router			/api/v1/films/{id} [get]
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
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmData", err, err.Error())
		return
	}

	logs.Logger.Debug("film:", film)
	logs.Logger.Debug("artists:", artists)

	userID, err := strconv.Atoi(gorilla_ctx.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmData", err, err.Error())
	}

	isSub, err := h.SubClient.CheckSub(
		context.Background(),
		&subscription.UserID{
			ID: strconv.Itoa(userID),
		},
	)

	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetFilmData.checkSub", err, err.Error())
	}

	if isSub.Status != true {
		film.MediaPath = ""
	}

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
//
//	@Summary		Get cast page
//	@Description	Get a list of films based on the cast name.
//	@Tags			Films
//	@Param			cast	path	string	true	"The Films of the Cast you want to retrieve."
//	@Produce		json
//
// @Success 		200 {json} []domain.Video
// @Failure 		400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
//
//	@Router			/api/v1/films/cast/{id} [get]
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
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
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

// GetTopRate godoc
//
//	@Summary		Get top rate information
//	@Description	Get information about the most rated film.
//	@Tags			Films
//	@Param			rate	path	string	true	"The top rate Film  you want to retrieve."
//	@Produce		json
//
// @Success 		200 {json} []domain.Video
// @Failure 		400 {json} domain.Response
// @Failure 		404 {json} domain.Response
// @Failure 		500 {json} domain.Response
//
//	@Router			/api/v1/films/top/rate [get]
func (h *FilmsHandler) GetTopRate(w http.ResponseWriter, r *http.Request) {
	film, err := h.FilmsUsecase.GetTopRate()
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "films_http", "GetTopRate", err, "Failed to get top")
		return
	}

	logs.Logger.Debug("films_http GetTopRate:", film)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"film": film,
		},
		http.StatusOK,
	)
}

func (h *FilmsHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films_http", "GetRecommendations", err, err.Error())
		return
	}

	rec, err := h.FilmsUsecase.GetRecommendations(userID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "films_http", "GetRecommendations", err, "Failed to get recommendations")
		return
	}

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"recommendations": rec,
		},
		http.StatusOK,
	)
}
