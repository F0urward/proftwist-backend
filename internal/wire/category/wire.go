//go:build wireinject
// +build wireinject

package category

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/metrics"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	loggingmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/logging"
	metricsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/metrics"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

func InitializeMetrics() metrics.Metrics {
	wire.Build(
		Metrics,
	)
	return nil
}

func InitializeCategoryHttpServer(cfg *config.Config, log logger.Logger, mtrs metrics.Metrics) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		CategorySet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
		metricsmiddleware.NewMetricsMiddleware,
	)
	return &httpServer.HttpServer{}
}
