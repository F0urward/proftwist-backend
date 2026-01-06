//go:build wireinject
// +build wireinject

package moderation

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

func InitializeModerationGrpcServer(cfg *config.Config) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		ModerationSet,
		AllGrpcRegistrars,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}
