package http

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func NewCsrfHandler(mainRouter *mux.Router) {

	mainRouter.HandleFunc("/api/v1/csrf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
	}).Methods(http.MethodGet, http.MethodOptions)

}
