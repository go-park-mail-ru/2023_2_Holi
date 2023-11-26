package http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"github.com/gorilla/context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FavouritesHandler struct {
	FavouritesUsecase domain.FavouritesUsecase
}

func NewFavouritesHandler(router *mux.Router, fu domain.FavouritesUsecase) {
	handler := &FavouritesHandler{
		FavouritesUsecase: fu,
	}

	router.HandleFunc("/v1/video/favourites/{id}", handler.AddToFavourites).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/video/favourites", handler.GetAllFavourites).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/video/favourites/{id}", handler.RemoveFromFavourites).Methods(http.MethodDelete, http.MethodOptions)
}

// AddToFavourites godoc
//
//	@Summary		Adds a video to favourites.
//	@Description	Adds a film or a whole series to favourites by id.
//	@Tags			Favourites
//	@Param			id	path	int	true	"The id of the video you want to add."
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/video/favourites/{id} [post]
func (h *FavouritesHandler) AddToFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "Add", err, err.Error())
		return
	}
	logs.Logger.Debug("AddToFavourites path param id: ", videoID)

	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "Add", err, err.Error())
	}
	logs.Logger.Debug("AddToFavourites user id: ", userID)

	err = h.FavouritesUsecase.AddToFavourites(videoID, userID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "Add", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveFromFavourites godoc
//
//	@Summary		Deletes a video from favourites.
//	@Description	Deletes a film or a whole series from favourites by id.
//	@Tags			Favourites
//	@Param			id	path	int	true	"The id of the video you want to delete."
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	object{err=string}
//	@Failure		404	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/video/favourites/{id} [delete]
func (h *FavouritesHandler) RemoveFromFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "RemoveFromFavourites", err, err.Error())
		return
	}
	logs.Logger.Debug("RemoveFromFavourites path param id: ", videoID)

	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "RemoveFromFavourites", err, err.Error())
	}
	logs.Logger.Debug("RemoveFromFavourites user id: ", userID)

	err = h.FavouritesUsecase.RemoveFromFavourites(videoID, userID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "RemoveFromFavourites", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllFavourites godoc
//
//	@Summary		Retrieves all video from favourites.
//	@Description	Retrieves all video from favourites.
//	@Tags			Favourites
//	@Produce		json
//	@Success		200	{object}	object{body=[]domain.VideoResponse}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/video/favourites/{id} [get]
func (h *FavouritesHandler) GetAllFavourites(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(context.Get(r, "userID").(string))
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "RemoveFromFavourites", err, err.Error())
	}
	logs.Logger.Debug("GetAllFavourites user id: ", userID)

	videos, err := h.FavouritesUsecase.GetAllFavourites(userID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "RemoveFromFavourites", err, err.Error())
		return
	}

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"videos": videos,
		},
		http.StatusOK,
	)
}
