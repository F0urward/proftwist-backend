package moderation

import (
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
)

func AllGrpcRegistrars(
	moderationGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		moderationGrpcRegistrar,
	}
}
