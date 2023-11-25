package http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/microcosm-cc/bluemonday"

	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	ProfileUsecase domain.ProfileUsecase
	Sanitizer      *bluemonday.Policy
}

func NewProfileHandler(router *mux.Router, pu domain.ProfileUsecase, s *bluemonday.Policy) {
	handler := &ProfileHandler{
		ProfileUsecase: pu,
		Sanitizer:      s,
	}

	router.HandleFunc("/v1/profile/{id}", handler.GetUserData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/profile/update", handler.UpdateProfile).Methods(http.MethodPost, http.MethodOptions)
}

// GetUserData godoc
//
//	@Summary		Get user by id
//	@Description	Get user data by id
//	@Tags			profile
//	@Param			id	path	int	true	"The user id you want to retrieve."
//	@Produce		json
//	@Success		200	{json}	domain.User
//	@Failure		400	{json}	ApiResponse
//	@Failure		404	{json}	ApiResponse
//	@Failure		500	{json}	ApiResponse
//	@Router			/api/v1/profile/{id} [get]
func (h *ProfileHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "GetUserData", err, err.Error())
		return
	}

	user, err := h.ProfileUsecase.GetUserData(userID)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "GetUserData", err, err.Error())
		return
	}

	su := domain.SanitizeUser(user, h.Sanitizer)
	logs.Logger.Debug("user:", su)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"user": su,
		},
		http.StatusOK,
	)
}

// UpdateProfile godoc
//
//	@Summary		update profile
//	@Description	update user data in db and return it
//	@Tags			profile
//	@Produce		json
//	@Accept			json
//	@Param			body	body		domain.UserRequest	true	"user that must be updated"
//	@Success		200		{object}	domain.Response{body=object{user=domain.User}}
//	@Failure		400		{json}		domain.Response
//	@Failure		403		{json}		domain.Response
//	@Failure		500		{json}		domain.Response
//	@Router			/api/v1/profile/update [post]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var newUser domain.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "UpdateProfile", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Need to update user for:", newUser)
	defer h.CloseAndAlert(r.Body)

	if len(newUser.ImageData) != 0 {
		newUser.ImagePath, err = h.ProfileUsecase.UploadImage(newUser.ID, newUser.ImageData)
		if err != nil {
			http.Error(w, `{"err":"`+err.Error()+`"}`, domain.GetHttpStatusCode(err))
			logs.LogError(logs.Logger, "profile_http", "UpdateProfile", err, "Failed to upload user image")
			return
		}
	}

	updatedUser, err := h.ProfileUsecase.UpdateUser(newUser)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
		logs.LogError(logs.Logger, "http", "UpdateProfile", err, err.Error())
		return
	}
	logs.Logger.Debug("Updated user:", updatedUser)

	su := domain.SanitizeUser(updatedUser, h.Sanitizer)
	logs.Logger.Debug("user:", su)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"user": su,
		},
		http.StatusOK,
	)
}

func (h *ProfileHandler) CloseAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logs.LogError(logs.Logger, "http", "CloseAndAlert", err, err.Error())
	}
}
