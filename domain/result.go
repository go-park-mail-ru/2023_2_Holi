package domain

import "net/http"

type Response struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

func CreateError(w http.ResponseWriter, errString string) {

}

func CreateResponse(w http.ResponseWriter, result map[string]interface{}) {

}
