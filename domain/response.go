package domain

import (
	"encoding/json"
	"net/http"
)

//easyjson:json
type Response struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

func WriteError(w http.ResponseWriter, errString string, status int) {
	//w.WriteHeader(status)
	//response := &Response{Err: errString}
	//easyjson.MarshalToHTTPResponseWriter(easyjson.Marshaler(response), w)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Response{Err: errString})
}

func WriteResponse(w http.ResponseWriter, result map[string]interface{}, status int) {
	//w.WriteHeader(status)
	//response := &Response{Body: result}
	//easyjson.MarshalToHTTPResponseWriter(easyjson.Marshaler(response), w)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Response{Body: result})
}
