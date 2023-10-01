package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logs"
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
	var credentials domain.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "Login", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Login credentials:", credentials)
	defer a.CloseAndAlert(r.Body)

	if credentials.Password == "" || credentials.Name == "" {
		http.Error(w, `{"err":"`+domain.ErrWrongCredentials.Error()+`"}`, http.StatusForbidden)
		logs.LogError(logs.Logger, "http", "Login", err, "Credentials are empy")
		return
	}

	session, err := a.AuthUsecase.Login(credentials)
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
		logs.LogError(logs.Logger, "http", "Login", err, "Failed to login")
		return
	}
	logs.Logger.Debug("Login: session:", session)

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
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusUnauthorized)
			logs.LogError(logs.Logger, "http", "Logout", err, "No cookie")
			return
		}

		http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "Logout", err, "Failed to get cookie")
		return
	}

	sessionToken := c.Value
	if sessionToken == "" {
		http.Error(w, `{"err":"`+domain.ErrUnauthorized.Error()+`"}`, http.StatusUnauthorized)
		logs.LogError(logs.Logger, "http", "Logout", err, "Session token is empty")
		return
	}
	logs.Logger.Debug("Logout: session token:", c)

	if err = a.AuthUsecase.Logout(sessionToken); err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusInternalServerError)
		logs.LogError(logs.Logger, "http", "Logout", err, "Failed to logout")
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
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusBadRequest)
		logs.LogError(logs.Logger, "http", "Register", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Register user:", user)
	defer a.CloseAndAlert(r.Body)

	if user.Name == "" || user.Password == "" {
		http.Error(w, `{"err":"name or password is empty"}`, http.StatusForbidden)
		logs.LogError(logs.Logger, "http", "Register", err, "User fields are empty")
		return
	}

	if id, err := a.AuthUsecase.Register(user); err != nil {
		http.Error(w, `{"err":"`+err.Error()+`"}`, getStatusCode(err))
		logs.LogError(logs.Logger, "http", "Register", err, "Failed to register")
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

func (a *AuthHandler) CloseAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logs.LogError(logs.Logger, "http", "CloseAndAlert", err, "Failed to close body")
	}
}
