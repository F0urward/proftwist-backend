package logging

import (
	"context"
	"time"

	"github.com/F0urward/proftwist-backend/internal/entities/domains"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/pkg/logger"
	"google.golang.org/grpc"
)

type LoggingUnaryServerInterceptor struct {
	log logger.Logger
}

func NewLoggingUnaryServerInterceptor(log logger.Logger) *LoggingUnaryServerInterceptor {
	return &LoggingUnaryServerInterceptor{
		log: log,
	}
}

func (l *LoggingUnaryServerInterceptor) LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		requestLogger := l.log.WithFields(map[string]interface{}{
			"grpc.method": info.FullMethod,
			"grpc.type":   "unary",
		})

		requestLogger.Info("gRPC request started")

		ctx = ctxutil.WithLogging(ctx, requestLogger)

		resp, err := handler(ctx, req)

		requestLogger.WithFields(map[string]interface{}{
			"took":  time.Since(start).String(),
			"error": err != nil,
		}).Info("gRPC request completed")

		return resp, err
	}
}

func WithLogging(ctx context.Context, log logger.Logger) context.Context {
	return context.WithValue(ctx, domains.LoggerKey{}, log)
}
