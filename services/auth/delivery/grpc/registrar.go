package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

type AuthGrpcRegistrar struct {
	server authclient.AuthServiceServer
}

func NewAuthGrpcRegistrar(server authclient.AuthServiceServer) grpcServer.GrpcRegistrar {
	return &AuthGrpcRegistrar{
		server: server,
	}
}

func (r *AuthGrpcRegistrar) RegisterServer(s *grpcServer.GrpcServer) {
	authclient.RegisterAuthServiceServer(s.Server, r.server)
}
