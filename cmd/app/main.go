package main

import (
	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/wire"
)

func main() {
	cfg := config.New()

	httpServer := wire.InitializeHttpServer(cfg)

	grpcServer := wire.InitializeGrpcServer(cfg)

	go httpServer.Run()

	grpcServer.Run()
}
