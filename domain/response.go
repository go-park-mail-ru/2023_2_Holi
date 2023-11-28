package domain

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

func WriteError(w http.ResponseWriter, errString string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Response{Err: errString})
}

func WriteResponse(w http.ResponseWriter, result map[string]interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Response{Body: result})
}
