//go:build wireinject
// +build wireinject

package auth

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
)

func InitializeAuthHttpServer(cfg *config.Config) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		AuthSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeAuthGrpcServer(cfg *config.Config) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		AuthSet,
		AllGrpcRegistrars,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}
