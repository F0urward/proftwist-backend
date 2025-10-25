package grpc

import "github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"

func (s *GrpcServer) RegisterServices() {
	roadmapclient.RegisterRoadmapServiceServer(s.Server, s.RoadmapServer)
}
