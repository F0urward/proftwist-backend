package roadmap

import (
	"github.com/F0urward/proftwist-backend/internal/metrics"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	aiHttp "github.com/F0urward/proftwist-backend/services/ai/delivery/http"
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
	aiHttpRegistrar *aiHttp.AIHttpRegistrar,
) []httpServer.HttpRegistrar {
	return []httpServer.HttpRegistrar{
		roadmapHttpRegistrar,
		aiHttpRegistrar,
	}
}

func AllGrpcRegistrars(
	roadmapGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		roadmapGrpcRegistrar,
	}
}
