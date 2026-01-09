//go:build wireinject
// +build wireinject

package moderation

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	logginginterceptor "github.com/F0urward/proftwist-backend/internal/server/interceptor/logging"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

func InitializeModerationGrpcServer(cfg *config.Config, log logger.Logger) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		ModerationSet,
		AllGrpcRegistrars,
		grpcServer.New,
		logginginterceptor.NewLoggingUnaryServerInterceptor,
	)
	return &grpcServer.GrpcServer{}
}
