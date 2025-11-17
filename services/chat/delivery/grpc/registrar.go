package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

type ChatGrpcRegistrar struct {
	server chatclient.ChatServiceServer
}

func NewChatGrpcRegistrar(server chatclient.ChatServiceServer) grpcServer.GrpcRegistrar {
	return &ChatGrpcRegistrar{
		server: server,
	}
}

func (r *ChatGrpcRegistrar) RegisterServer(s *grpcServer.GrpcServer) {
	chatclient.RegisterChatServiceServer(s.Server, r.server)
}
