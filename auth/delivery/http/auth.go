package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"
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

	router.HandleFunc("/api/v1/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/api/v1/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/api/v1/auth/logout", handler.Logout).Methods("POST")
}

// Login godoc
// @Summary      login user
// @Description  create user session and put it into cookie
// @Tags         auth
// @Accept       json
// @Success      204
// @Failure      400  {string} string "{"error":"<error message>"}"
// @Failure      403  {string} string "{"error":"<error message>"}"
// @Failure      404  {string} string "{"error":"<error message>"}"
// @Failure      500  {string} string "{"error":"<error message>"}"
// @Router       /api/v1/auth/login [post]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	auth, _ := a.auth(r)
	if auth == true {
		http.Error(w, `{"err":"you must be unauthorised"}`, http.StatusForbidden)
	}

	var credentials domain.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	defer CloseAndAlert(r.Body)

	if credentials.Password == "" || credentials.Email == "" {
		http.Error(w, `{"err":"`+domain.ErrWrongCredentials.Error()+`"}`, http.StatusForbidden)
		return
	}

	credentials.Email = strings.TrimSpace(credentials.Email)

	if err = checkCredentials(credentials); err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
	}

	session, err := a.AuthUsecase.Login(credentials)
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusNoContent)
}

// Logout godoc
// @Summary      logout user
// @Description  delete current session and nullify cookie
// @Tags         auth
// @Success      204
// @Failure      400  {string} string "{"error":"<error message>"}"
// @Failure      403  {string} string "{"error":"<error message>"}"
// @Failure      404  {string} string "{"error":"<error message>"}"
// @Failure      500  {string} string "{"error":"<error message>"}"
// @Router       /api/v1/auth/logout [post]
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth, err := a.auth(r)
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
		return
	}
	if !auth {
		http.Error(w, `{"err":"`+domain.ErrUnauthorized.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	c, err := r.Cookie("session_token")
	sessionToken := c.Value
	if err = a.AuthUsecase.Logout(sessionToken); err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	w.WriteHeader(http.StatusNoContent)
}

// Register godoc
// @Summary      register user
// @Description  add new user to db and return it id
// @Tags         auth
// @Produce      json
// @Accept       json
// @Success      200  {object} Result
// @Failure      400  {string} string "{"error":"<error message>"}"
// @Failure      500  {string} string "{"error":"<error message>"}"
// @Router       /api/v1/auth/register [post]
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	auth, _ := a.auth(r)
	if auth == true {
		http.Error(w, `{"err":"you must be unauthorised"}`, http.StatusForbidden)
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	defer CloseAndAlert(r.Body)

	user.Email = strings.TrimSpace(user.Email)
	if err = checkCredentials(domain.Credentials{Email: user.Email, Password: user.Password}); err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
	}

	var id int
	if id, err = a.AuthUsecase.Register(user); err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
		return
	}

	session, err := a.AuthUsecase.Login(domain.Credentials{Email: user.Email, Password: user.Password})
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	body := map[string]interface{}{
		"id": id,
	}
	json.NewEncoder(w).Encode(&Result{Body: body})
}

func (a *AuthHandler) auth(r *http.Request) (bool, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, domain.ErrUnauthorized
		}

		return false, domain.ErrBadRequest
	}
	if c.Expires.After(time.Now()) {
		return false, domain.ErrUnauthorized
	}
	sessionToken := c.Value
	exists, err := a.AuthUsecase.IsAuth(sessionToken)
	if err != nil {
		return false, domain.ErrInternalServerError
	}
	if !exists {
		return false, domain.ErrUnauthorized
	}

	return true, nil
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
	case domain.ErrAlreadyExists:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func CloseAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println("err: ", err)
	}
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func checkCredentials(cred domain.Credentials) error {
	if cred.Email == "" || cred.Password == "" {
		return domain.ErrWrongCredentials
	}

	if !valid(cred.Email) {
		return domain.ErrWrongCredentials
	}

	return nil
}
