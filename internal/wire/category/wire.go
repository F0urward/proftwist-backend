//go:build wireinject
// +build wireinject

package category

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
)

func InitializeCategoryHttpServer(cfg *config.Config) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		CategorySet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
	)
	return &httpServer.HttpServer{}
}
