//go:build wireinject
// +build wireinject

package roadmap

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

func InitializeRoadmapHttpServer(cfg *config.Config, log logger.Logger, mtrs metrics.Metrics) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		RoadmapSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
		metricsmiddleware.NewMetricsMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeRoadmapGrpcServer(cfg *config.Config, log logger.Logger, mtrs metrics.Metrics) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		RoadmapSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
		metricsinterceptor.NewMetricsUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
