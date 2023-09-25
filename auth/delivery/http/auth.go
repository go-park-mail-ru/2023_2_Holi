package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"2023_2_Holi/domain"
)

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandler(router *mux.Router, u domain.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: u,
	}

	router.HandleFunc("api/v1/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("api/v1/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("api/v1/auth/logout", handler.Logout).Methods("POST")
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials domain.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	defer closeAndAlert(r.Body)

	if credentials.Password == "" || credentials.Name == "" {
		http.Error(w, `{"error":"`+domain.ErrWrongCredentials.Error()+`"}`, http.StatusForbidden)
		return
	}

	session, err := a.AuthUsecase.Login(credentials)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, getStatusCode(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Expires: session.ExpiresAt,
	})
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	if err = a.AuthUsecase.Logout(sessionToken); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	defer closeAndAlert(r.Body)

	if id, err := a.AuthUsecase.Register(user); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, getStatusCode(err))
		return

	} else {
		body := map[string]interface{}{
			"id": id,
		}
		json.NewEncoder(w).Encode(&Result{Body: body})
	}
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
	case domain.ErrWrongCredentials:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func closeAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
