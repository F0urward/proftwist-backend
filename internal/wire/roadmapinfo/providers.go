package roadmapinfo

import (
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
)

func AllHttpRegistrars(
	roadmapInfoHttpRegistrar httpServer.HttpRegistrar,
) []httpServer.HttpRegistrar {
	return []httpServer.HttpRegistrar{
		roadmapInfoHttpRegistrar,
	}
}

func AllGrpcRegistrars(
	roadmapInfoGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		roadmapInfoGrpcRegistrar,
	}
}
