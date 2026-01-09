package ctxutil

import (
	"context"

	"github.com/F0urward/proftwist-backend/internal/entities/domains"
	"github.com/F0urward/proftwist-backend/pkg/logger"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func WithLogging(ctx context.Context, logger logger.Logger) context.Context {
	return context.WithValue(ctx, domains.LoggerKey{}, logger)
}

func GetLogger(ctx context.Context) logger.Logger {
	if logger, ok := ctx.Value(domains.LoggerKey{}).(logger.Logger); ok {
		return logger
	}

	return GetDefaultLogger()
}

func GetDefaultLogger() logger.Logger {
	return logrus.NewLogrusLogger()
}
