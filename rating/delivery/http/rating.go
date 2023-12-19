package http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"encoding/json"
	"github.com/gorilla/context"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type RatingHandler struct {
	RatingUsecase domain.RatingUsecase
}

func NewRatingHandler(router *mux.Router, fu domain.RatingUsecase) {
	handler := &RatingHandler{
		RatingUsecase: fu,
	}

	router.HandleFunc("/v1/video/rating", handler.AddRate).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/video/rating/{id}", handler.DeleteRate).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/v1/video/rating/check/{id}", handler.Rated).Methods(http.MethodGet, http.MethodOptions)
}

// AddRate godoc
//
//	@Summary		Adds the rate to the video.
//	@Description	Adds the rate to the video.
//	@Tags			Rating
//	@Param			body	body		domain.Rate	true	"user credentials"
//	@Success		200	{object}	object{body=object{rating=float}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/video/rating [post]
func (h *RatingHandler) AddRate(w http.ResponseWriter, r *http.Request) {
	var rate domain.Rate
	err := json.NewDecoder(r.Body).Decode(&rate)
	if err != nil {
		domain.WriteError(w, "Invalid body", http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "AddRate", err, "Failed to decode json from body")
		return
	}
	if rate.VideoID <= 0 || rate.Rate <= 0 {
		domain.WriteError(w, "Invalid body", http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "AddRate", domain.ErrBadRequest, "Invalid rate or video id")
		return
	}
	logs.Logger.Debug("AddRate rate:", rate)
	defer h.CloseAndAlert(r.Body)

	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "AddRate", err, err.Error())
	}
	logs.Logger.Debug("AddRate user id: ", userID)
	rate.UserID = userID

	rating, err := h.RatingUsecase.Add(rate)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "Add", err, err.Error())
		return
	}

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"rating": rating,
		},
		http.StatusOK,
	)
}

// DeleteRate godoc
//
//	@Summary		Deletes the rate.
//	@Description	Deletes the user rate for the video.
//	@Tags			Rating
//	@Param			id	path	int	true	"The id of the video."
//	@Success		200	{object}	object{body=object{rating=float}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		404	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/video/rating/{id} [delete]
func (h *RatingHandler) DeleteRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "DeleteRate", err, err.Error())
		return
	}
	if videoID <= 0 {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "DeleteRate", domain.ErrBadRequest, "invalid video id")
		return
	}
	logs.Logger.Debug("DeleteRate path param id: ", videoID)

	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "DeleteRate", err, err.Error())
	}
	logs.Logger.Debug("DeleteRate user id: ", userID)

	rating, err := h.RatingUsecase.Remove(domain.Rate{UserID: userID, VideoID: videoID})
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "DeleteRate", err, err.Error())
		return
	}

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"rating": rating,
		},
		http.StatusOK,
	)
}

// Rated godoc
//
//	@Summary		checks is rated
//	@Description	checks if the video was rated by the user
//	@Tags			Rating
//	@Param			id	path	int	true	"The id of the video."
//	@Success		200	{object}	object{body=object{isRated=bool}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/v1/video/rating/check/{id} [post]
func (h *RatingHandler) Rated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "Rated", err, err.Error())
		return
	}
	if videoID <= 0 {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "Rated", domain.ErrBadRequest, "invalid video id")
		return
	}
	logs.Logger.Debug("Rated path param id: ", videoID)

	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "Rated", err, err.Error())
	}
	logs.Logger.Debug("Rated user id: ", userID)

	rated, err := h.RatingUsecase.Rated(domain.Rate{UserID: userID, VideoID: videoID})
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "Rated", err, err.Error())
		return
	}

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"isRated": rated,
		},
		http.StatusOK,
	)
}

func (h *RatingHandler) CloseAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logs.LogError(logs.Logger, "http", "CloseAndAlert", err, "Failed to close body")
	}
}
