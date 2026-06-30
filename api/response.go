package main

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Data json.RawMessage `json:"data"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorData `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data json.RawMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{Data: data})
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorData{
			Code:    code,
			Message: msg,
		},
	})
}
