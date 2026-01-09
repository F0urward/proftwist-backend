//go:build wireinject
// +build wireinject

package category

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	loggingmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/logging"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

func InitializeCategoryHttpServer(cfg *config.Config, log logger.Logger) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		CategorySet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
		loggingmiddleware.NewLoggingMiddleware,
	)
	return &httpServer.HttpServer{}
}
