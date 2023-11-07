package http

import (
	"2023_2_Holi/domain"
	"net/http"
	"time"

	logs "2023_2_Holi/logger"

	"github.com/gorilla/mux"
)

type CsrfHandler struct {
	Token *domain.HashToken
}

func NewCsrfHandler(mainRouter *mux.Router, t *domain.HashToken) {
	handler := &CsrfHandler{
		Token: t,
	}
	mainRouter.HandleFunc("/api/v1/csrf", func(w http.ResponseWriter, r *http.Request) {
		token, err := handler.Token.Create("uuid.NewString()", time.Now().Add(24*time.Hour).Unix())
		if err != nil {
			http.Error(w, `{"err":"`+err.Error()+`"}`, domain.GetStatusCode(err))
			logs.LogError(logs.Logger, "csrf token", "creation error:", err, "Failed to create")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "csrf-token",
			Value: token,
			Path:  "/",
		})
	}).Methods(http.MethodGet, http.MethodOptions)

}
