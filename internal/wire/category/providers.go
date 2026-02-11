package category

import (
	"github.com/F0urward/proftwist-backend/internal/metrics"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() metrics.Metrics {
	reg := prometheus.NewRegistry()

	wrapped := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"service": "proftwist-category-service",
		},
		reg,
	)

	return metrics.NewMetrics(reg, wrapped)
}

func AllHttpRegistrars(
	categoryHttpRegistrar httpServer.HttpRegistrar,
) []httpServer.HttpRegistrar {
	return []httpServer.HttpRegistrar{
		categoryHttpRegistrar,
	}
}
