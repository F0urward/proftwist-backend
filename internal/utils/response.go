package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}

func JSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")

	rawBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to marshal response: %v", err)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(rawBytes)
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

func JSONError(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(ctx, w, statusCode, NewErrorResponse(message))
}
