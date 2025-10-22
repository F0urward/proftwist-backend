package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
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
	const op = "utils.JSONResponse"
	logger := logctx.GetLogger(context.Background()).WithField("op", op)

	w.Header().Set("Content-Type", "application/json")

	rawBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("failed to marshal response")
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(rawBytes)
	if err != nil {
		logger.WithError(err).Error("failed to write response")
		return
	}
}

func JSONError(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(ctx, w, statusCode, NewErrorResponse(message))
}
