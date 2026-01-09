//go:build wireinject
// +build wireinject

package friend

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

func InitializeFriendHttpServer(cfg *config.Config, log logger.Logger) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		FriendSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeFriendGrpcServer(cfg *config.Config, log logger.Logger) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		FriendSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
