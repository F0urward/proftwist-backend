package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/moderationclient"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

type ModerationGrpcRegistrar struct {
	server moderationclient.ModerationServiceServer
}

func NewModerationGrpcRegistrar(server moderationclient.ModerationServiceServer) grpcServer.GrpcRegistrar {
	return &ModerationGrpcRegistrar{
		server: server,
	}
}

func (r *ModerationGrpcRegistrar) RegisterServer(s *grpcServer.GrpcServer) {
	moderationclient.RegisterModerationServiceServer(s.Server, r.server)
}
