package roadmapinfo

import (
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
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
