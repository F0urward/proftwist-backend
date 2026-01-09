//go:build wireinject
// +build wireinject

package roadmap

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

func InitializeRoadmapHttpServer(cfg *config.Config, log logger.Logger) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		RoadmapSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeRoadmapGrpcServer(cfg *config.Config, log logger.Logger) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		RoadmapSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
