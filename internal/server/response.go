package httpserver

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, msg string) {
	msgWithErr := errorResponse{Message: msg}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if msg != "" {
		_ = json.NewEncoder(w).Encode(msgWithErr)
	}
}
