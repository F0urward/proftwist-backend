package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapinfoclient"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

type RoadmapInfoGrpcRegistrar struct {
	server roadmapinfoclient.RoadmapInfoServiceServer
}

func NewRoadmapInfoGrpcRegistrar(server roadmapinfoclient.RoadmapInfoServiceServer) grpcServer.GrpcRegistrar {
	return &RoadmapInfoGrpcRegistrar{
		server: server,
	}
}

func (r *RoadmapInfoGrpcRegistrar) RegisterServer(s *grpcServer.GrpcServer) {
	roadmapinfoclient.RegisterRoadmapInfoServiceServer(s.Server, r.server)
}
