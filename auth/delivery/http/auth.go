package http

import (
	"2023_2_Holi/domain"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewAuthHandler(u domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}

	http.HandleFunc("/auth/login", handler.Login)
	http.HandleFunc("/auth/register", handler.Register)
	http.HandleFunc("/auth/logout", handler.Logout)
}

func (a *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	session, err := a.UserUsecase.Login(user)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
	}

	http.SetCookie(w, &http.Cookie{
		Name:    session.SessionData,
		Value:   session.Token,
		Expires: session.ExpiresAt,
	})
}

func (a *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (a *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {

}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
