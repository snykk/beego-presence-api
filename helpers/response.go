package helpers

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Status  bool        `json:"status"`           // Indicates success or failure
	Message string      `json:"message"`          // Response message
	Data    interface{} `json:"data,omitempty"`   // Actual response data
	Error   interface{} `json:"errors,omitempty"` // Error details, if any
}

func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err != nil {
		switch e := err.(type) {
		case error:
			err = e.Error()
		default:
			// No action required; keep err as is
		}
	}

	json.NewEncoder(w).Encode(BaseResponse{
		Status:  false,
		Message: message,
		Error:   err,
	})
}
