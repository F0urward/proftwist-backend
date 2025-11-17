package roadmap

import (
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
)

func AllHttpRegistrars(
	roadmapHttpRegistrar httpServer.HttpRegistrar,
) []httpServer.HttpRegistrar {
	return []httpServer.HttpRegistrar{
		roadmapHttpRegistrar,
	}
}

func AllGrpcRegistrars(
	roadmapGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		roadmapGrpcRegistrar,
	}
}
