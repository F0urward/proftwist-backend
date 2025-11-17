//go:build wireinject
// +build wireinject

package roadmapinfo

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
)

func InitializeRoadmapInfoHttpServer(cfg *config.Config) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		RoadmapInfoSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeRoadmapInfoGrpcServer(cfg *config.Config) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		RoadmapInfoSet,
		AllGrpcRegistrars,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}
