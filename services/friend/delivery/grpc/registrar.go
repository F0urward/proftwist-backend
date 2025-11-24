package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/friendclient"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

type FriendGrpcRegistrar struct {
	server friendclient.FriendServiceServer
}

func NewFriendGrpcRegistrar(server friendclient.FriendServiceServer) grpcServer.GrpcRegistrar {
	return &FriendGrpcRegistrar{
		server: server,
	}
}

func (r *FriendGrpcRegistrar) RegisterServer(s *grpcServer.GrpcServer) {
	friendclient.RegisterFriendServiceServer(s.Server, r.server)
}
