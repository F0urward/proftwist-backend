package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
)

func (s *GrpcServer) RegisterServices() {
	roadmapclient.RegisterRoadmapServiceServer(s.Server, s.RoadmapServer)
	roadmapinfoclient.RegisterRoadmapInfoServiceServer(s.Server, s.RoadmapInfoServer)
	authclient.RegisterAuthServiceServer(s.Server, s.AuthServer)
	chatclient.RegisterChatServiceServer(s.Server, s.ChatServer)
}
