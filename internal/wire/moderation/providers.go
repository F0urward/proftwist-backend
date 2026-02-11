package moderation

import (
	"github.com/F0urward/proftwist-backend/internal/metrics"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() metrics.Metrics {
	reg := prometheus.NewRegistry()

	wrapped := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"service": "proftwist-moderation-service",
		},
		reg,
	)

	return metrics.NewMetrics(reg, wrapped)
}

func AllGrpcRegistrars(
	moderationGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		moderationGrpcRegistrar,
	}
}
