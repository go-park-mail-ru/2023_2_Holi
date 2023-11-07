package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewCsrfHandler(mainRouter *mux.Router) {

	mainRouter.HandleFunc("/api/v1/csrf", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet, http.MethodOptions)

}
