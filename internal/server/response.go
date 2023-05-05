package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
