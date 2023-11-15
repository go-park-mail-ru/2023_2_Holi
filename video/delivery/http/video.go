package http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VideoHandler struct {
	VideoUsecase domain.VideoUsecase
}

func NewVideoHandler(router *mux.Router, fu domain.VideoUsecase) {
	handler := &VideoHandler{
		VideoUsecase: fu,
	}

	router.HandleFunc("/v1/video/favourites/{id}", handler.AddToFavourites).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/video/favourites/{id}", handler.).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/video/favourites/{id}", handler.DeleteFromFavourites).Methods(http.MethodDelete, http.MethodOptions)
}

// AddToFavourites godoc
// @Summary 		Adds a video to favourites.
// @Description 	Adds a films or a whole series to favourites by id
// @Tags 			Favourites
// @Param 			id path int true "The id of the video you want to add."
// @Produce 		json
// @Success         204
// @Failure         400  {object} object{err=string}
// @Failure         404  {object} object{err=string}
// @Failure         500  {object} object{err=string}
// @Router 			/api/v1/video/favourites/{id} [post]
func (h *VideoHandler) AddToFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "AddToFavourites", err, err.Error())
		return
	}

	err := h.VideoUsecase.AddToFavourites(id)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "http", "AddToFavourites", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
