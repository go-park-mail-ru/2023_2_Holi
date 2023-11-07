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
	json.NewEncoder(w).Encode(&Response{Err: errString})
	w.WriteHeader(status)
}

func WriteResponse(w http.ResponseWriter, result map[string]interface{}, status int) {
	json.NewEncoder(w).Encode(&Response{Body: result})
	w.WriteHeader(status)
}
