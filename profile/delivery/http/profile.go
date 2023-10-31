package profile_http

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

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

	router.HandleFunc("/v1/profile/{id}", handler.GetUserData).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/v1/profile/update", handler.UpdateProfile).Methods(http.MethodPost, http.MethodOptions)
}

// GetUserData godoc
// @Summary 		Get user by id
// @Description 	Get user data by id
// @Tags 			profile
// @Param 			id path int true "The user id you want to retrieve."
// @Produce 		json
// @Success 		200 {json} domain.User
// @Failure			400 {json} ApiResponse
// @Failure 		404 {json} ApiResponse
// @Failure 		500 {json} ApiResponse
// @Router 			/api/v1/profile/{id} [get]
func (h *ProfileHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "profile_http", "GetUserData", err, "failed to read user id")
		return
	}

	user, err := h.ProfileUsecase.GetUserData(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "profile_http", "GetUserData", err, err.Error())
		return
	}

	response := ApiResponse{
		Status: http.StatusOK,
		Body: map[string]interface{}{
			"user": user,
		},
	}

	logs.Logger.Debug("user:", user)
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile godoc
// @Summary      update profile
// @Description  update user data in db and return it
// @Tags         profile
// @Produce      json
// @Accept       json
// @Success      200  {json} Result
// @Failure      400  {json} Result
// @Failure      403  {json} Result
// @Failure      500  {json} Result
// @Router       /api/v1/profile/update [post]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var newUser domain.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "profile_http", "UpdateProfile", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Need to update user for:", newUser)
	defer h.CloseAndAlert(r.Body)

	if len(newUser.ImageData) != 0 {
		newUser.ImagePath, err = h.ProfileUsecase.UploadImage(newUser.ID, newUser.ImageData)
	}

	updatedUser, err := h.ProfileUsecase.UpdateUser(newUser)
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "profile_http", "UpdateProfile", err, "Failed to update user profile")
		return
	}
	logs.Logger.Debug("Updated user:", updatedUser)

	response := ApiResponse{
		Status: http.StatusOK,
		Body: map[string]interface{}{
			"user": updatedUser,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (h *ProfileHandler) CloseAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logs.LogError(logs.Logger, "auth_http", "CloseAndAlert", err, "Failed to close body")
	}
}
