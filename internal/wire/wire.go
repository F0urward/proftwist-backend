//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	"github.com/F0urward/proftwist-backend/internal/wire/initializers"
)

func InitializeHttpServer(cfg *config.Config) *httpServer.HttpServer {
	wire.Build(
		initializers.HTTPServerSet,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeGrpcServer(cfg *config.Config) *grpcServer.GrpcServer {
	wire.Build(
		initializers.GRPCServerSet,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}
