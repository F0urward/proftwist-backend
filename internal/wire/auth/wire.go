//go:build wireinject
// +build wireinject

package auth

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/metrics"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	logginginterceptor "github.com/F0urward/proftwist-backend/internal/server/interceptor/logging"
	metricsinterceptor "github.com/F0urward/proftwist-backend/internal/server/interceptor/metrics"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	loggingmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/logging"
	metricsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/metrics"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

func InitializeMetrics() metrics.Metrics {
	wire.Build(
		Metrics,
	)
	return nil
}

func InitializeAuthHttpServer(cfg *config.Config, log logger.Logger, mtrs metrics.Metrics) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		AuthSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
		metricsmiddleware.NewMetricsMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeAuthGrpcServer(cfg *config.Config, log logger.Logger, metrics metrics.Metrics) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		AuthSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
		metricsinterceptor.NewMetricsUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
