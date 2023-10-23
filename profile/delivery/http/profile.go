package profile_http

import (
	"2023_2_Holi/domain"
	"net/http"

	"github.com/gorilla/mux"
)

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type ProfileHandler struct {
	ProfileUsecase domain.ProfileUsecase
}

func NewProfileHandler(router *mux.Router, pu domain.ProfileUsecase) {
	handler := &ProfileHandler{
		ProfileUsecase: pu,
	}

	router.HandleFunc("/api/v1/profile/{id}", handler.GetProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/update", handler.UpdateProfile).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/image", handler.UploadImage).Methods(http.MethodPost, http.MethodOptions)
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

}

func (h *ProfileHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

}
