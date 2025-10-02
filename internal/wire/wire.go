//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
)

func InitializeHttpServer(cfg *config.Config) *httpServer.HttpServer {
	wire.Build(
		httpServer.New,
	)
	return &httpServer.HttpServer{}
}

