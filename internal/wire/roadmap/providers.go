package roadmap

import (
	"github.com/F0urward/proftwist-backend/internal/metrics"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() metrics.Metrics {
	reg := prometheus.NewRegistry()

	wrapped := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"service": "proftwist-roadmap-service",
		},
		reg,
	)

	return metrics.NewMetrics(reg, wrapped)
}

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
