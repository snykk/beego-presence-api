package helpers

import (
	"encoding/json"
	"net/http"
)

type baseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := baseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := baseResponse{
		Status:  false,
		Message: message,
		Error:   err.Error(),
	}
	json.NewEncoder(w).Encode(response)
}
