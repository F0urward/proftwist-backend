//go:build wireinject
// +build wireinject

package roadmapinfo

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	logginginterceptor "github.com/F0urward/proftwist-backend/internal/server/interceptor/logging"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	loggingmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/logging"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

func InitializeRoadmapInfoHttpServer(cfg *config.Config, log logger.Logger) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		RoadmapInfoSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeRoadmapInfoGrpcServer(cfg *config.Config, log logger.Logger) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		RoadmapInfoSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
