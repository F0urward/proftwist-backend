package main

import (
	"github.com/F0urward/proftwist-backend/config"
	authWire "github.com/F0urward/proftwist-backend/internal/wire/auth"
)

func main() {
	cfg := config.New()

	httpServer := authWire.InitializeAuthHttpServer(cfg)

	grpcServer := authWire.InitializeAuthGrpcServer(cfg)

	go httpServer.Run()

	grpcServer.Run()
}
