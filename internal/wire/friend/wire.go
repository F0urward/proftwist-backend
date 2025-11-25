//go:build wireinject
// +build wireinject

package friend

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
)

func InitializeFriendHttpServer(cfg *config.Config) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		FriendSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeFriendGrpcServer(cfg *config.Config) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		FriendSet,
		AllGrpcRegistrars,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}
