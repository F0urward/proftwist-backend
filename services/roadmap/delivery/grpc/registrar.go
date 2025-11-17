package grpc

import (
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/roadmapclient"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

type RoadmapGrpcRegistrar struct {
	server roadmapclient.RoadmapServiceServer
}

func NewRoadmapGrpcRegistrar(server roadmapclient.RoadmapServiceServer) grpcServer.GrpcRegistrar {
	return &RoadmapGrpcRegistrar{
		server: server,
	}
}

func (r *RoadmapGrpcRegistrar) RegisterServer(s *grpcServer.GrpcServer) {
	roadmapclient.RegisterRoadmapServiceServer(s.Server, r.server)
}
