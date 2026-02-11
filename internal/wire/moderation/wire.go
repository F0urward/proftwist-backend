//go:build wireinject
// +build wireinject

package moderation

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/metrics"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	logginginterceptor "github.com/F0urward/proftwist-backend/internal/server/interceptor/logging"
	metricsinterceptor "github.com/F0urward/proftwist-backend/internal/server/interceptor/metrics"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

func InitializeMetrics() metrics.Metrics {
	wire.Build(
		Metrics,
	)
	return nil
}

func InitializeModerationGrpcServer(cfg *config.Config, log logger.Logger, mtrs metrics.Metrics) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		ModerationSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
		metricsinterceptor.NewMetricsUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
