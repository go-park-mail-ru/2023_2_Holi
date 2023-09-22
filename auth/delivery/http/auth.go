package http

import (
	"2023_2_Holi/domain"
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

}

func (a *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (a *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {

}
