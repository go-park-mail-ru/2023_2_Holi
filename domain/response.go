package domain

import (
	"github.com/mailru/easyjson"
	"net/http"
)

//easyjson:json
type Response struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

func WriteError(w http.ResponseWriter, errString string, status int) {
	w.WriteHeader(status)
	response := &Response{Err: errString}
	easyjson.MarshalToHTTPResponseWriter(easyjson.Marshaler(response), w)
}

func WriteResponse(w http.ResponseWriter, result map[string]interface{}, status int) {
	w.WriteHeader(status)
	response := &Response{Body: result}
	easyjson.MarshalToHTTPResponseWriter(easyjson.Marshaler(response), w)
}
